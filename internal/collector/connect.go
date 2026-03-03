package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type ConnectCollector struct{}

func NewConnectCollector() *ConnectCollector { return &ConnectCollector{} }

func (c *ConnectCollector) Name() string { return "connect" }

func (c *ConnectCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := connect.NewFromConfig(cfg)
	paginator := connect.NewListInstancesPaginator(client, &connect.ListInstancesInput{})

	counts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, inst := range page.InstanceSummaryList {
			status := string(inst.InstanceStatus)
			if status == "" {
				status = "UNKNOWN"
			}
			counts[status]++
		}
	}

	for status, count := range counts {
		metrics.ConnectInstances.WithLabelValues(account, accountName, region, status).Set(count)
	}
	return nil
}
