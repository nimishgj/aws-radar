package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type NeptuneCollector struct{}

func NewNeptuneCollector() *NeptuneCollector {
	return &NeptuneCollector{}
}

func (c *NeptuneCollector) Name() string {
	return "neptune"
}

func (c *NeptuneCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := neptune.NewFromConfig(cfg)
	paginator := neptune.NewDescribeDBClustersPaginator(client, &neptune.DescribeDBClustersInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.DBClusters))
	}

	metrics.NeptuneClusters.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
