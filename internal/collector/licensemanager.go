package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/licensemanager"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type LicenseManagerCollector struct{}

func NewLicenseManagerCollector() *LicenseManagerCollector {
	return &LicenseManagerCollector{}
}

func (c *LicenseManagerCollector) Name() string {
	return "licensemanager"
}

func (c *LicenseManagerCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := licensemanager.NewFromConfig(cfg)

	var count float64
	var nextToken *string
	for {
		output, err := client.ListLicenseConfigurations(ctx, &licensemanager.ListLicenseConfigurationsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return err
		}
		count += float64(len(output.LicenseConfigurations))
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}
		nextToken = output.NextToken
	}

	metrics.LicenseManagerConfigurations.WithLabelValues(account, accountName, region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("license_configurations", count).
		Msg("License Manager collection completed")

	return nil
}
