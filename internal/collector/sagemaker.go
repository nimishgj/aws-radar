package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type SageMakerCollector struct{}

func NewSageMakerCollector() *SageMakerCollector { return &SageMakerCollector{} }

func (c *SageMakerCollector) Name() string { return "sagemaker" }

func (c *SageMakerCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := sagemaker.NewFromConfig(cfg)
	paginator := sagemaker.NewListEndpointsPaginator(client, &sagemaker.ListEndpointsInput{})

	counts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, ep := range page.Endpoints {
			status := string(ep.EndpointStatus)
			if status == "" {
				status = "UNKNOWN"
			}
			counts[status]++
		}
	}
	for status, count := range counts {
		metrics.SageMakerEndpoints.WithLabelValues(account, accountName, region, status).Set(count)
	}
	return nil
}
