package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type AppStreamCollector struct{}

func NewAppStreamCollector() *AppStreamCollector { return &AppStreamCollector{} }

func (c *AppStreamCollector) Name() string { return "appstream" }

func (c *AppStreamCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := appstream.NewFromConfig(cfg)

	counts := make(map[string]float64)
	var nextToken *string
	for {
		page, err := client.DescribeFleets(ctx, &appstream.DescribeFleetsInput{NextToken: nextToken})
		if err != nil {
			return err
		}
		for _, fleet := range page.Fleets {
			state := string(fleet.State)
			if state == "" {
				state = "UNKNOWN"
			}
			counts[state]++
		}
		if page.NextToken == nil || *page.NextToken == "" {
			break
		}
		nextToken = page.NextToken
	}
	for state, count := range counts {
		metrics.AppStreamFleets.WithLabelValues(account, accountName, region, state).Set(count)
	}
	return nil
}
