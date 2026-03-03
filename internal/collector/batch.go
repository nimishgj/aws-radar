package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type BatchCollector struct{}

func NewBatchCollector() *BatchCollector {
	return &BatchCollector{}
}

func (c *BatchCollector) Name() string {
	return "batch"
}

func (c *BatchCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := batch.NewFromConfig(cfg)
	paginator := batch.NewDescribeJobQueuesPaginator(client, &batch.DescribeJobQueuesInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.JobQueues))
	}

	metrics.BatchJobQueues.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
