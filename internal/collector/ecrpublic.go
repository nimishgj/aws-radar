package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type ECRPublicCollector struct{}

func NewECRPublicCollector() *ECRPublicCollector { return &ECRPublicCollector{} }

func (c *ECRPublicCollector) Name() string { return "ecrpublic" }

func (c *ECRPublicCollector) Collect(ctx context.Context, cfg aws.Config, account, accountName string) error {
	cfg.Region = "us-east-1"
	client := ecrpublic.NewFromConfig(cfg)
	paginator := ecrpublic.NewDescribeRepositoriesPaginator(client, &ecrpublic.DescribeRepositoriesInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Repositories))
	}

	metrics.ECRPublicRepositories.WithLabelValues(account, accountName).Set(count)
	return nil
}
