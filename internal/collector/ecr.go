package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type ECRCollector struct{}

func NewECRCollector() *ECRCollector {
	return &ECRCollector{}
}

func (c *ECRCollector) Name() string {
	return "ecr"
}

func (c *ECRCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := ecr.NewFromConfig(cfg)

	var count float64
	paginator := ecr.NewDescribeRepositoriesPaginator(client, &ecr.DescribeRepositoriesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Repositories))
	}

	metrics.ECRRepositories.WithLabelValues(account, accountName, region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("repository_count", count).
		Msg("ECR repository collection completed")

	return nil
}
