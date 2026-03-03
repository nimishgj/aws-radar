package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53resolver"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type Route53ResolverCollector struct{}

func NewRoute53ResolverCollector() *Route53ResolverCollector { return &Route53ResolverCollector{} }

func (c *Route53ResolverCollector) Name() string { return "route53resolver" }

func (c *Route53ResolverCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := route53resolver.NewFromConfig(cfg)

	endpointCounts := make(map[string]float64)
	epPaginator := route53resolver.NewListResolverEndpointsPaginator(client, &route53resolver.ListResolverEndpointsInput{})
	for epPaginator.HasMorePages() {
		page, err := epPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, ep := range page.ResolverEndpoints {
			direction := string(ep.Direction)
			if direction == "" {
				direction = "UNKNOWN"
			}
			status := string(ep.Status)
			if status == "" {
				status = "UNKNOWN"
			}
			endpointCounts[direction+"|"+status]++
		}
	}
	for key, count := range endpointCounts {
		parts := splitKey(key, 2)
		metrics.Route53ResolverEndpoints.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	ruleCounts := make(map[string]float64)
	rulePaginator := route53resolver.NewListResolverRulesPaginator(client, &route53resolver.ListResolverRulesInput{})
	for rulePaginator.HasMorePages() {
		page, err := rulePaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, rule := range page.ResolverRules {
			ruleType := string(rule.RuleType)
			if ruleType == "" {
				ruleType = "UNKNOWN"
			}
			ruleCounts[ruleType]++
		}
	}
	for ruleType, count := range ruleCounts {
		metrics.Route53ResolverRules.WithLabelValues(account, accountName, region, ruleType).Set(count)
	}

	return nil
}
