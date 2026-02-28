package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type EFSCollector struct{}

func NewEFSCollector() *EFSCollector {
	return &EFSCollector{}
}

func (c *EFSCollector) Name() string {
	return "efs"
}

func (c *EFSCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := efs.NewFromConfig(cfg)

	var count float64
	paginator := efs.NewDescribeFileSystemsPaginator(client, &efs.DescribeFileSystemsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.FileSystems))
	}

	metrics.EFSFileSystems.WithLabelValues(account, accountName, region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("filesystem_count", count).
		Msg("EFS filesystem collection completed")

	return nil
}
