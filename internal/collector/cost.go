package collector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	cetypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type CostCollector struct {
	frequency string

	mu          sync.RWMutex
	lastFetch   time.Time
	periodStart string
	total       float64
	byService   map[string]float64
}

func NewCostCollector(frequency string) *CostCollector {
	freq := strings.ToLower(strings.TrimSpace(frequency))
	switch freq {
	case "hourly", "daily", "weekly", "monthly":
	case "":
		freq = "daily"
	default:
		freq = "daily"
	}
	return &CostCollector{
		frequency: freq,
		byService: make(map[string]float64),
	}
}

func (c *CostCollector) Name() string {
	return "cost_explorer"
}

func (c *CostCollector) Collect(ctx context.Context, cfg aws.Config, account, accountName string) error {
	now := time.Now().UTC()
	if c.shouldRefresh(now) {
		if err := c.refresh(ctx, cfg, now); err != nil {
			return err
		}
	}

	c.mu.RLock()
	periodStart := c.periodStart
	total := c.total
	byService := make(map[string]float64, len(c.byService))
	for service, amount := range c.byService {
		byService[service] = amount
	}
	c.mu.RUnlock()

	if periodStart == "" {
		return nil
	}

	metrics.CostTotal.WithLabelValues(account, accountName, periodStart).Set(total)
	for service, amount := range byService {
		metrics.CostByService.WithLabelValues(account, accountName, service, periodStart).Set(amount)
	}

	return nil
}

func (c *CostCollector) shouldRefresh(now time.Time) bool {
	c.mu.RLock()
	lastFetch := c.lastFetch
	frequency := c.frequency
	c.mu.RUnlock()

	if lastFetch.IsZero() {
		return true
	}

	switch frequency {
	case "hourly":
		return now.Format("2006-01-02T15") != lastFetch.UTC().Format("2006-01-02T15")
	case "daily":
		return now.Format("2006-01-02") != lastFetch.UTC().Format("2006-01-02")
	case "weekly":
		nowYear, nowWeek := now.ISOWeek()
		lastYear, lastWeek := lastFetch.UTC().ISOWeek()
		return nowYear != lastYear || nowWeek != lastWeek
	case "monthly":
		return now.Format("2006-01") != lastFetch.UTC().Format("2006-01")
	default:
		// Unknown values fall back to daily behavior.
		return now.Format("2006-01-02") != lastFetch.UTC().Format("2006-01-02")
	}
}

func (c *CostCollector) refresh(ctx context.Context, cfg aws.Config, now time.Time) error {
	ceCfg := cfg.Copy()
	if ceCfg.Region == "" {
		ceCfg.Region = "us-east-1"
	}
	client := costexplorer.NewFromConfig(ceCfg)

	startTime, endTime, granularity := c.queryWindow(now)
	start := startTime.Format("2006-01-02")
	end := endTime.Format("2006-01-02")

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &cetypes.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
		Granularity: granularity,
		Metrics:     []string{"UnblendedCost"},
		GroupBy: []cetypes.GroupDefinition{
			{
				Type: cetypes.GroupDefinitionTypeDimension,
				Key:  aws.String("SERVICE"),
			},
		},
	}

	costByService := make(map[string]float64)
	var periodStart string
	var total float64
	var nextToken *string

	for {
		input.NextPageToken = nextToken
		out, err := client.GetCostAndUsage(ctx, input)
		if err != nil {
			return fmt.Errorf("get cost and usage: %w", err)
		}

		for _, result := range out.ResultsByTime {
			if periodStart == "" {
				periodStart = aws.ToString(result.TimePeriod.Start)
			}
			for _, group := range result.Groups {
				if len(group.Keys) == 0 {
					continue
				}
				service := group.Keys[0]
				metric, ok := group.Metrics["UnblendedCost"]
				if !ok {
					continue
				}
				amount, err := strconv.ParseFloat(aws.ToString(metric.Amount), 64)
				if err != nil {
					return fmt.Errorf("parse cost amount for service %q: %w", service, err)
				}
				costByService[service] += amount
				total += amount
			}
		}

		nextToken = out.NextPageToken
		if nextToken == nil || *nextToken == "" {
			break
		}
	}

	c.mu.Lock()
	c.byService = costByService
	c.total = total
	c.periodStart = periodStart
	c.lastFetch = now
	c.mu.Unlock()

	return nil
}

func (c *CostCollector) queryWindow(now time.Time) (time.Time, time.Time, cetypes.Granularity) {
	now = now.UTC()
	switch c.frequency {
	case "hourly":
		// Cost Explorer GetCostAndUsage accepts dates, so hourly still queries current UTC day.
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return dayStart, dayStart.AddDate(0, 0, 1), cetypes.GranularityDaily
	case "weekly":
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return dayStart.AddDate(0, 0, -7), dayStart, cetypes.GranularityDaily
	case "monthly":
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		return monthStart, now, cetypes.GranularityDaily
	case "daily":
		fallthrough
	default:
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return dayStart.AddDate(0, 0, -1), dayStart, cetypes.GranularityDaily
	}
}
