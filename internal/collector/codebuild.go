package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type CodeBuildCollector struct{}

func NewCodeBuildCollector() *CodeBuildCollector {
	return &CodeBuildCollector{}
}

func (c *CodeBuildCollector) Name() string {
	return "codebuild"
}

func (c *CodeBuildCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := codebuild.NewFromConfig(cfg)
	paginator := codebuild.NewListProjectsPaginator(client, &codebuild.ListProjectsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Projects))
	}

	metrics.CodeBuildProjects.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
