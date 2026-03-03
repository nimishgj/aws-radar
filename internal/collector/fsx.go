package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/fsx"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type FSxCollector struct{}

func NewFSxCollector() *FSxCollector {
	return &FSxCollector{}
}

func (c *FSxCollector) Name() string {
	return "fsx"
}

func (c *FSxCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := fsx.NewFromConfig(cfg)
	paginator := fsx.NewDescribeFileSystemsPaginator(client, &fsx.DescribeFileSystemsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.FileSystems))
	}

	metrics.FSxFileSystems.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
