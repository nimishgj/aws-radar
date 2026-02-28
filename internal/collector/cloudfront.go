package collector

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type CloudFrontCollector struct{}

func NewCloudFrontCollector() *CloudFrontCollector {
	return &CloudFrontCollector{}
}

func (c *CloudFrontCollector) Name() string {
	return "cloudfront"
}

func (c *CloudFrontCollector) Collect(ctx context.Context, cfg aws.Config, account, accountName string) error {
	// CloudFront is a global service, use us-east-1
	cfg.Region = "us-east-1"
	client := cloudfront.NewFromConfig(cfg)

	counts := make(map[string]float64)

	paginator := cloudfront.NewListDistributionsPaginator(client, &cloudfront.ListDistributionsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		if page.DistributionList == nil || page.DistributionList.Items == nil {
			continue
		}

		for _, dist := range page.DistributionList.Items {
			priceClass := string(dist.PriceClass)
			enabled := strconv.FormatBool(aws.ToBool(dist.Enabled))

			key := priceClass + "|" + enabled
			counts[key]++
		}
	}

	// Update metrics
	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.CloudFrontDistributions.WithLabelValues(
			account,
			accountName,
			parts[0], // price_class
			parts[1], // enabled
		).Set(count)
	}

	log.Debug().
		Int("distribution_combinations", len(counts)).
		Msg("CloudFront collection completed")

	return nil
}
