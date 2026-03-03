package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type DMSCollector struct{}

func NewDMSCollector() *DMSCollector { return &DMSCollector{} }

func (c *DMSCollector) Name() string { return "dms" }

func (c *DMSCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := databasemigrationservice.NewFromConfig(cfg)
	paginator := databasemigrationservice.NewDescribeReplicationInstancesPaginator(client, &databasemigrationservice.DescribeReplicationInstancesInput{})

	counts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, inst := range page.ReplicationInstances {
			status := aws.ToString(inst.ReplicationInstanceStatus)
			if status == "" {
				status = "unknown"
			}
			counts[status]++
		}
	}
	for status, count := range counts {
		metrics.DMSReplicationInstances.WithLabelValues(account, accountName, region, status).Set(count)
	}
	return nil
}
