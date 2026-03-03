package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/memorydb"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type MemoryDBCollector struct{}

func NewMemoryDBCollector() *MemoryDBCollector {
	return &MemoryDBCollector{}
}

func (c *MemoryDBCollector) Name() string {
	return "memorydb"
}

func (c *MemoryDBCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := memorydb.NewFromConfig(cfg)
	paginator := memorydb.NewDescribeClustersPaginator(client, &memorydb.DescribeClustersInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Clusters))
	}

	metrics.MemoryDBClusters.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
