package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type EMRCollector struct{}

func NewEMRCollector() *EMRCollector {
	return &EMRCollector{}
}

func (c *EMRCollector) Name() string {
	return "emr"
}

func (c *EMRCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := emr.NewFromConfig(cfg)
	paginator := emr.NewListClustersPaginator(client, &emr.ListClustersInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Clusters))
	}

	metrics.EMRClusters.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
