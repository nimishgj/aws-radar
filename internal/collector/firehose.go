package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type FirehoseCollector struct{}

func NewFirehoseCollector() *FirehoseCollector {
	return &FirehoseCollector{}
}

func (c *FirehoseCollector) Name() string {
	return "firehose"
}

func (c *FirehoseCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := firehose.NewFromConfig(cfg)
	var count float64
	var startName *string
	for {
		page, err := client.ListDeliveryStreams(ctx, &firehose.ListDeliveryStreamsInput{
			ExclusiveStartDeliveryStreamName: startName,
		})
		if err != nil {
			return err
		}
		count += float64(len(page.DeliveryStreamNames))
		if page.HasMoreDeliveryStreams == nil || !*page.HasMoreDeliveryStreams || len(page.DeliveryStreamNames) == 0 {
			break
		}
		last := page.DeliveryStreamNames[len(page.DeliveryStreamNames)-1]
		startName = aws.String(last)
	}

	metrics.FirehoseDeliveryStreams.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
