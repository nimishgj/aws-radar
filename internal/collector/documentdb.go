package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type DocumentDBCollector struct{}

func NewDocumentDBCollector() *DocumentDBCollector {
	return &DocumentDBCollector{}
}

func (c *DocumentDBCollector) Name() string {
	return "documentdb"
}

func (c *DocumentDBCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := docdb.NewFromConfig(cfg)
	paginator := docdb.NewDescribeDBClustersPaginator(client, &docdb.DescribeDBClustersInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.DBClusters))
	}

	metrics.DocumentDBClusters.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
