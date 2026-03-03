package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type KinesisCollector struct{}

func NewKinesisCollector() *KinesisCollector {
	return &KinesisCollector{}
}

func (c *KinesisCollector) Name() string {
	return "kinesis"
}

func (c *KinesisCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := kinesis.NewFromConfig(cfg)
	paginator := kinesis.NewListStreamsPaginator(client, &kinesis.ListStreamsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.StreamNames))
	}

	metrics.KinesisStreams.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
