package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type APIGatewayV2Collector struct{}

func NewAPIGatewayV2Collector() *APIGatewayV2Collector {
	return &APIGatewayV2Collector{}
}

func (c *APIGatewayV2Collector) Name() string {
	return "apigatewayv2"
}

func (c *APIGatewayV2Collector) Collect(ctx context.Context, cfg aws.Config, region string) error {
	client := apigatewayv2.NewFromConfig(cfg)

	counts := make(map[string]float64)
	var nextToken *string
	for {
		page, err := client.GetApis(ctx, &apigatewayv2.GetApisInput{
			NextToken: nextToken,
		})
		if err != nil {
			return err
		}
		for _, api := range page.Items {
			protocol := string(api.ProtocolType)
			if protocol == "" {
				protocol = "unknown"
			}
			counts[protocol]++
		}
		if page.NextToken == nil || len(*page.NextToken) == 0 {
			break
		}
		nextToken = page.NextToken
	}

	for protocol, count := range counts {
		metrics.APIGatewayV2APIs.WithLabelValues(region, protocol).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("protocols", len(counts)).
		Msg("API Gateway v2 API collection completed")

	return nil
}
