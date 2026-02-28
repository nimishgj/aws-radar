package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type OpenSearchCollector struct{}

func NewOpenSearchCollector() *OpenSearchCollector {
	return &OpenSearchCollector{}
}

func (c *OpenSearchCollector) Name() string {
	return "opensearch"
}

func (c *OpenSearchCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := opensearch.NewFromConfig(cfg)

	output, err := client.ListDomainNames(ctx, &opensearch.ListDomainNamesInput{})
	if err != nil {
		return err
	}

	count := float64(len(output.DomainNames))
	metrics.OpenSearchDomains.WithLabelValues(account, accountName, region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("domain_count", count).
		Msg("OpenSearch domain collection completed")

	return nil
}
