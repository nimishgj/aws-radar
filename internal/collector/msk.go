package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type MSKCollector struct{}

func NewMSKCollector() *MSKCollector {
	return &MSKCollector{}
}

func (c *MSKCollector) Name() string {
	return "msk"
}

func (c *MSKCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := kafka.NewFromConfig(cfg)
	paginator := kafka.NewListClustersV2Paginator(client, &kafka.ListClustersV2Input{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.ClusterInfoList))
	}

	metrics.MSKClusters.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
