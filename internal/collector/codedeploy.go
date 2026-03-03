package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type CodeDeployCollector struct{}

func NewCodeDeployCollector() *CodeDeployCollector {
	return &CodeDeployCollector{}
}

func (c *CodeDeployCollector) Name() string {
	return "codedeploy"
}

func (c *CodeDeployCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := codedeploy.NewFromConfig(cfg)
	paginator := codedeploy.NewListApplicationsPaginator(client, &codedeploy.ListApplicationsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Applications))
	}

	metrics.CodeDeployApplications.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
