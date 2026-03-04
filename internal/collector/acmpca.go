package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type ACMPCACollector struct{}

func NewACMPCACollector() *ACMPCACollector {
	return &ACMPCACollector{}
}

func (c *ACMPCACollector) Name() string {
	return "acmpca"
}

func (c *ACMPCACollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := acmpca.NewFromConfig(cfg)
	counts := make(map[string]float64)

	paginator := acmpca.NewListCertificateAuthoritiesPaginator(client, &acmpca.ListCertificateAuthoritiesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, ca := range page.CertificateAuthorities {
			status := string(ca.Status)
			if status == "" {
				status = "unknown"
			}
			caType := string(ca.Type)
			if caType == "" {
				caType = "unknown"
			}
			key := status + "|" + caType
			counts[key]++
		}
	}

	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.ACMPCACertificateAuthorities.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("ca_combinations", len(counts)).
		Msg("ACM PCA collection completed")

	return nil
}
