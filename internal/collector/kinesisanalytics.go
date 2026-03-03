package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type KinesisAnalyticsCollector struct{}

func NewKinesisAnalyticsCollector() *KinesisAnalyticsCollector {
	return &KinesisAnalyticsCollector{}
}

func (c *KinesisAnalyticsCollector) Name() string {
	return "kinesisanalytics"
}

func (c *KinesisAnalyticsCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := kinesisanalyticsv2.NewFromConfig(cfg)
	paginator := kinesisanalyticsv2.NewListApplicationsPaginator(client, &kinesisanalyticsv2.ListApplicationsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.ApplicationSummaries))
	}

	metrics.KinesisAnalyticsApplications.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
