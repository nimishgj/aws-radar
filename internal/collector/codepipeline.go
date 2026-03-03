package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type CodePipelineCollector struct{}

func NewCodePipelineCollector() *CodePipelineCollector {
	return &CodePipelineCollector{}
}

func (c *CodePipelineCollector) Name() string {
	return "codepipeline"
}

func (c *CodePipelineCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := codepipeline.NewFromConfig(cfg)
	paginator := codepipeline.NewListPipelinesPaginator(client, &codepipeline.ListPipelinesInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Pipelines))
	}

	metrics.CodePipelinePipelines.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
