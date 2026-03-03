package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type CloudFormationCollector struct{}

func NewCloudFormationCollector() *CloudFormationCollector {
	return &CloudFormationCollector{}
}

func (c *CloudFormationCollector) Name() string {
	return "cloudformation"
}

func (c *CloudFormationCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := cloudformation.NewFromConfig(cfg)
	paginator := cloudformation.NewListStacksPaginator(client, &cloudformation.ListStacksInput{})

	var count float64
	statusCounts := make(map[types.StackStatus]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, summary := range page.StackSummaries {
			count++
			statusCounts[summary.StackStatus]++
		}
	}

	metrics.CloudFormationStacks.WithLabelValues(account, accountName, region).Set(count)
	for status, statusCount := range statusCounts {
		metrics.CloudFormationStacksByStatus.WithLabelValues(account, accountName, region, string(status)).Set(statusCount)
	}
	return nil
}
