package collector

import (
	"context"
	"strings"

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
	listenerCounts := make(map[string]float64)
	targetGroupCounts := make(map[string]float64)
	rulesPerALB := make(map[string]float64)
	azCountPerLB := make(map[string]float64)
	subnetCountPerLB := make(map[string]float64)
	lbMetaByArn := make(map[string][3]string) // arn -> [name, type, scheme]

	paginator := elbv2.NewDescribeLoadBalancersPaginator(client, &elbv2.DescribeLoadBalancersInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, lb := range page.LoadBalancers {
			lbType := string(lb.Type)
			if lbType == "" {
				lbType = "unknown"
			}
			scheme := string(lb.Scheme)
			if scheme == "" {
				scheme = "unknown"
			}
			ipType := string(lb.IpAddressType)
			if ipType == "" {
				ipType = "unknown"
			}
			state := string(lb.State.Code)
			if state == "" {
				state = "unknown"
			}
			lbName := aws.ToString(lb.LoadBalancerName)
			if lbName == "" {
				arnParts := strings.Split(aws.ToString(lb.LoadBalancerArn), "/")
				lbName = arnParts[len(arnParts)-1]
			}
			lbArn := aws.ToString(lb.LoadBalancerArn)
			lbMetaByArn[lbArn] = [3]string{lbName, lbType, scheme}

			key := lbType + "|" + scheme
			counts[key]++

			detailedKey := lbType + "|" + scheme + "|" + ipType + "|" + state
			detailedCounts[detailedKey]++

			azSet := make(map[string]struct{})
			subnetSet := make(map[string]struct{})
			for _, az := range lb.AvailabilityZones {
				azName := aws.ToString(az.ZoneName)
				if azName != "" {
					azSet[azName] = struct{}{}
				}
				subnetID := aws.ToString(az.SubnetId)
				if subnetID != "" {
					subnetSet[subnetID] = struct{}{}
				}
			}
			lbKey := lbName + "|" + lbType + "|" + scheme
			azCountPerLB[lbKey] = float64(len(azSet))
			subnetCountPerLB[lbKey] = float64(len(subnetSet))

			listenerPaginator := elbv2.NewDescribeListenersPaginator(client, &elbv2.DescribeListenersInput{
				LoadBalancerArn: lb.LoadBalancerArn,
			})
			for listenerPaginator.HasMorePages() {
				listenerPage, listenerErr := listenerPaginator.NextPage(ctx)
				if listenerErr != nil {
					log.Warn().Err(listenerErr).Str("region", region).Str("load_balancer", lbName).Msg("Failed to describe ELBv2 listeners")
					break
				}
				for _, listener := range listenerPage.Listeners {
					protocol := string(listener.Protocol)
					if protocol == "" {
						protocol = "unknown"
					}
					listenerKey := lbType + "|" + scheme + "|" + protocol
					listenerCounts[listenerKey]++

					if lbType == "application" {
						rulePaginator := elbv2.NewDescribeRulesPaginator(client, &elbv2.DescribeRulesInput{
							ListenerArn: listener.ListenerArn,
						})
						for rulePaginator.HasMorePages() {
							rulePage, ruleErr := rulePaginator.NextPage(ctx)
							if ruleErr != nil {
								log.Warn().Err(ruleErr).Str("region", region).Str("load_balancer", lbName).Msg("Failed to describe ELBv2 listener rules")
								break
							}
							rulesPerALB[lbName] += float64(len(rulePage.Rules))
						}
					}
				}
			}
		}
	}

	targetGroupPaginator := elbv2.NewDescribeTargetGroupsPaginator(client, &elbv2.DescribeTargetGroupsInput{})
	for targetGroupPaginator.HasMorePages() {
		page, err := targetGroupPaginator.NextPage(ctx)
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to describe ELBv2 target groups")
			break
		}
		for _, tg := range page.TargetGroups {
			targetType := string(tg.TargetType)
			if targetType == "" {
				targetType = "unknown"
			}
			lbType := "unknown"
			if len(tg.LoadBalancerArns) > 0 {
				if meta, ok := lbMetaByArn[tg.LoadBalancerArns[0]]; ok {
					lbType = meta[1]
				}
			}
			key := lbType + "|" + targetType
			targetGroupCounts[key]++
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

	for key, count := range listenerCounts {
		parts := splitKey(key, 3)
		metrics.ELBV2Listeners.WithLabelValues(account, accountName, region, parts[0], parts[1], parts[2]).Set(count)
	}

	for key, count := range targetGroupCounts {
		parts := splitKey(key, 2)
		metrics.ELBV2TargetGroups.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	for lbName, count := range rulesPerALB {
		metrics.ELBV2RulesPerALB.WithLabelValues(account, accountName, region, lbName).Set(count)
	}

	for key, count := range azCountPerLB {
		parts := splitKey(key, 3)
		metrics.ELBV2AvailabilityZonesPerLB.WithLabelValues(account, accountName, region, parts[0], parts[1], parts[2]).Set(count)
	}

	for key, count := range subnetCountPerLB {
		parts := splitKey(key, 3)
		metrics.ELBV2SubnetsPerLB.WithLabelValues(account, accountName, region, parts[0], parts[1], parts[2]).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("elbv2_combinations", len(counts)).
		Int("elbv2_detailed_combinations", len(detailedCounts)).
		Int("elbv2_listener_combinations", len(listenerCounts)).
		Int("elbv2_target_group_combinations", len(targetGroupCounts)).
		Int("elbv2_rule_lb_count", len(rulesPerALB)).
		Msg("ELBv2 collection completed")

	return nil
}
