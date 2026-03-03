package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type KinesisVideoCollector struct{}

func NewKinesisVideoCollector() *KinesisVideoCollector { return &KinesisVideoCollector{} }

func (c *KinesisVideoCollector) Name() string { return "kinesisvideo" }

func (c *KinesisVideoCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := kinesisvideo.NewFromConfig(cfg)
	paginator := kinesisvideo.NewListStreamsPaginator(client, &kinesisvideo.ListStreamsInput{})
	counts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, stream := range page.StreamInfoList {
			status := string(stream.Status)
			if status == "" {
				status = "UNKNOWN"
			}
			counts[status]++
		}
	}
	for status, count := range counts {
		metrics.KinesisVideoStreams.WithLabelValues(account, accountName, region, status).Set(count)
	}
	return nil
}
