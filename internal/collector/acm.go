package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type ACMCollector struct{}

func NewACMCollector() *ACMCollector {
	return &ACMCollector{}
}

func (c *ACMCollector) Name() string {
	return "acm"
}

func (c *ACMCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := acm.NewFromConfig(cfg)

	counts := make(map[string]float64)

	paginator := acm.NewListCertificatesPaginator(client, &acm.ListCertificatesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, cert := range page.CertificateSummaryList {
			// Get certificate details
			descOutput, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
				CertificateArn: cert.CertificateArn,
			})
			if err != nil {
				log.Warn().
					Err(err).
					Str("certificate", aws.ToString(cert.CertificateArn)).
					Msg("Failed to describe certificate")
				continue
			}

			status := string(descOutput.Certificate.Status)
			certType := string(descOutput.Certificate.Type)

			key := status + "|" + certType
			counts[key]++
		}
	}

	// Update metrics
	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.ACMCertificates.WithLabelValues(account, accountName, region,
			parts[0], // status
			parts[1], // type
		).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("certificate_combinations", len(counts)).
		Msg("ACM collection completed")

	return nil
}
