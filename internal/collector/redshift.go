package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type RedshiftCollector struct{}

func NewRedshiftCollector() *RedshiftCollector {
	return &RedshiftCollector{}
}

func (c *RedshiftCollector) Name() string {
	return "redshift"
}

func (c *RedshiftCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := redshift.NewFromConfig(cfg)
	paginator := redshift.NewDescribeClustersPaginator(client, &redshift.DescribeClustersInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Clusters))
	}

	metrics.RedshiftClusters.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
