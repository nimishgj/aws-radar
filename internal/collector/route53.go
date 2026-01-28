package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type Route53Collector struct{}

func NewRoute53Collector() *Route53Collector {
	return &Route53Collector{}
}

func (c *Route53Collector) Name() string {
	return "route53"
}

func (c *Route53Collector) Collect(ctx context.Context, cfg aws.Config) error {
	// Route53 is a global service
	cfg.Region = "us-east-1"
	client := route53.NewFromConfig(cfg)

	var count float64 = 0

	paginator := route53.NewListHostedZonesPaginator(client, &route53.ListHostedZonesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		count += float64(len(page.HostedZones))
	}

	metrics.Route53HostedZones.WithLabelValues().Set(count)

	log.Debug().
		Float64("hosted_zone_count", count).
		Msg("Route53 collection completed")

	return nil
}
