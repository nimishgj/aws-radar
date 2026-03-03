package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type CloudTrailCollector struct{}

func NewCloudTrailCollector() *CloudTrailCollector {
	return &CloudTrailCollector{}
}

func (c *CloudTrailCollector) Name() string {
	return "cloudtrail"
}

func (c *CloudTrailCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := cloudtrail.NewFromConfig(cfg)
	output, err := client.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{
		IncludeShadowTrails: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	metrics.CloudTrailTrails.WithLabelValues(account, accountName, region).Set(float64(len(output.TrailList)))

	paginator := cloudtrail.NewListEventDataStoresPaginator(client, &cloudtrail.ListEventDataStoresInput{})
	statusCounts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, store := range page.EventDataStores {
			status := string(store.Status)
			if status == "" {
				status = "UNKNOWN"
			}
			statusCounts[status]++
		}
	}
	for status, count := range statusCounts {
		metrics.CloudTrailLakeEventDataStores.WithLabelValues(account, accountName, region, status).Set(count)
	}
	return nil
}
