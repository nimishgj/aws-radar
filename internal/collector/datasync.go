package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/datasync"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type DataSyncCollector struct{}

func NewDataSyncCollector() *DataSyncCollector { return &DataSyncCollector{} }

func (c *DataSyncCollector) Name() string { return "datasync" }

func (c *DataSyncCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := datasync.NewFromConfig(cfg)
	paginator := datasync.NewListTasksPaginator(client, &datasync.ListTasksInput{})
	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Tasks))
	}
	metrics.DataSyncTasks.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
