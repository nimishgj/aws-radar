package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type ELBCollector struct{}

func NewELBCollector() *ELBCollector {
	return &ELBCollector{}
}

func (c *ELBCollector) Name() string {
	return "elb"
}

func (c *ELBCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	// Collect Classic Load Balancers
	if err := c.collectClassic(ctx, cfg, region, account, accountName); err != nil {
		log.Warn().Err(err).Str("region", region).Msg("Failed to collect Classic ELB")
	}

	// Collect ALB/NLB
	if err := c.collectV2(ctx, cfg, region, account, accountName); err != nil {
		log.Warn().Err(err).Str("region", region).Msg("Failed to collect ELBv2")
	}

	return nil
}

func (c *ELBCollector) collectClassic(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := elb.NewFromConfig(cfg)

	counts := make(map[string]float64)

	output, err := client.DescribeLoadBalancers(ctx, &elb.DescribeLoadBalancersInput{})
	if err != nil {
		return err
	}

	for _, lb := range output.LoadBalancerDescriptions {
		scheme := aws.ToString(lb.Scheme)
		if scheme == "" {
			scheme = "internet-facing"
		}
		counts[scheme]++
	}

	for scheme, count := range counts {
		metrics.ELBClassic.WithLabelValues(account, accountName, region, scheme).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("classic_elb_count", len(output.LoadBalancerDescriptions)).
		Msg("Classic ELB collection completed")

	return nil
}

func (c *ELBCollector) collectV2(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := elbv2.NewFromConfig(cfg)

	counts := make(map[string]float64)
	detailedCounts := make(map[string]float64)

	paginator := elbv2.NewDescribeLoadBalancersPaginator(client, &elbv2.DescribeLoadBalancersInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, lb := range page.LoadBalancers {
			lbType := string(lb.Type)
			scheme := string(lb.Scheme)
			ipType := string(lb.IpAddressType)
			state := string(lb.State.Code)

			key := lbType + "|" + scheme
			counts[key]++

			detailedKey := lbType + "|" + scheme + "|" + ipType + "|" + state
			detailedCounts[detailedKey]++
		}
	}

	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.ELBV2.WithLabelValues(account, accountName, region,
			parts[0], // type
			parts[1], // scheme
		).Set(count)
	}

	for key, count := range detailedCounts {
		parts := splitKey(key, 4)
		metrics.ELBV2Detailed.WithLabelValues(account, accountName, region,
			parts[0], // type
			parts[1], // scheme
			parts[2], // ip_address_type
			parts[3], // state
		).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("elbv2_combinations", len(counts)).
		Int("elbv2_detailed_combinations", len(detailedCounts)).
		Msg("ELBv2 collection completed")

	return nil
}
