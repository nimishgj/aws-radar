package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type APIGatewayCollector struct{}

func NewAPIGatewayCollector() *APIGatewayCollector {
	return &APIGatewayCollector{}
}

func (c *APIGatewayCollector) Name() string {
	return "apigateway"
}

func (c *APIGatewayCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := apigateway.NewFromConfig(cfg)

	var count float64
	paginator := apigateway.NewGetRestApisPaginator(client, &apigateway.GetRestApisInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Items))
	}

	metrics.APIGatewayRestAPIs.WithLabelValues(account, accountName, region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("rest_api_count", count).
		Msg("API Gateway REST API collection completed")

	return nil
}
