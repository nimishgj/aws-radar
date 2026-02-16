package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type SSMCollector struct{}

func NewSSMCollector() *SSMCollector {
	return &SSMCollector{}
}

func (c *SSMCollector) Name() string {
	return "ssm"
}

func (c *SSMCollector) Collect(ctx context.Context, cfg aws.Config, region string) error {
	client := ssm.NewFromConfig(cfg)

	counts := make(map[string]float64)
	paginator := ssm.NewDescribeParametersPaginator(client, &ssm.DescribeParametersInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, param := range page.Parameters {
			paramType := string(param.Type)
			if paramType == "" {
				paramType = "unknown"
			}
			counts[paramType]++
		}
	}

	for paramType, count := range counts {
		metrics.SSMParameters.WithLabelValues(region, paramType).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("parameter_types", len(counts)).
		Msg("SSM parameter collection completed")

	return nil
}
