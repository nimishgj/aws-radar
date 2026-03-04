package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/fms"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type FMSCollector struct{}

func NewFMSCollector() *FMSCollector {
	return &FMSCollector{}
}

func (c *FMSCollector) Name() string {
	return "fms"
}

func (c *FMSCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := fms.NewFromConfig(cfg)
	counts := make(map[string]float64)

	var nextToken *string
	for {
		output, err := client.ListPolicies(ctx, &fms.ListPoliciesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return err
		}

		for _, policy := range output.PolicyList {
			resourceType := aws.ToString(policy.ResourceType)
			if resourceType == "" {
				resourceType = "unknown"
			}
			counts[resourceType]++
		}

		if output.NextToken == nil || *output.NextToken == "" {
			break
		}
		nextToken = output.NextToken
	}

	for resourceType, count := range counts {
		metrics.FMSPolicies.WithLabelValues(account, accountName, region, resourceType).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("resource_types", len(counts)).
		Msg("Firewall Manager collection completed")

	return nil
}
