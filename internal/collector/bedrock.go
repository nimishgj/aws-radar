package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type BedrockCollector struct{}

func NewBedrockCollector() *BedrockCollector { return &BedrockCollector{} }

func (c *BedrockCollector) Name() string { return "bedrock" }

func (c *BedrockCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := bedrock.NewFromConfig(cfg)
	paginator := bedrock.NewListCustomModelsPaginator(client, &bedrock.ListCustomModelsInput{})

	counts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, model := range page.ModelSummaries {
			status := string(model.ModelStatus)
			if status == "" {
				status = "UNKNOWN"
			}
			counts[status]++
		}
	}
	for status, count := range counts {
		metrics.BedrockCustomModels.WithLabelValues(account, accountName, region, status).Set(count)
	}
	return nil
}
