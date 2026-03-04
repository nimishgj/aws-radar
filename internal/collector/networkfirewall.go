package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	networkfirewall "github.com/aws/aws-sdk-go-v2/service/networkfirewall"
	nfTypes "github.com/aws/aws-sdk-go-v2/service/networkfirewall/types"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type NetworkFirewallCollector struct{}

func NewNetworkFirewallCollector() *NetworkFirewallCollector {
	return &NetworkFirewallCollector{}
}

func (c *NetworkFirewallCollector) Name() string {
	return "networkfirewall"
}

func (c *NetworkFirewallCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := networkfirewall.NewFromConfig(cfg)

	// List Firewalls
	var firewallCount float64
	var fwNextToken *string
	for {
		output, err := client.ListFirewalls(ctx, &networkfirewall.ListFirewallsInput{
			NextToken: fwNextToken,
		})
		if err != nil {
			return err
		}
		firewallCount += float64(len(output.Firewalls))
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}
		fwNextToken = output.NextToken
	}
	metrics.NetworkFirewallFirewalls.WithLabelValues(account, accountName, region).Set(firewallCount)

	// List Firewall Policies
	var policyCount float64
	var policyNextToken *string
	for {
		output, err := client.ListFirewallPolicies(ctx, &networkfirewall.ListFirewallPoliciesInput{
			NextToken: policyNextToken,
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to list Network Firewall policies")
			break
		}
		policyCount += float64(len(output.FirewallPolicies))
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}
		policyNextToken = output.NextToken
	}
	metrics.NetworkFirewallPolicies.WithLabelValues(account, accountName, region).Set(policyCount)

	// List Rule Groups (account-owned)
	var ruleGroupCount float64
	var rgNextToken *string
	for {
		output, err := client.ListRuleGroups(ctx, &networkfirewall.ListRuleGroupsInput{
			NextToken: rgNextToken,
			Scope:     nfTypes.ResourceManagedStatusAccount,
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to list Network Firewall rule groups")
			break
		}
		ruleGroupCount += float64(len(output.RuleGroups))
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}
		rgNextToken = output.NextToken
	}
	metrics.NetworkFirewallRuleGroups.WithLabelValues(account, accountName, region, "account").Set(ruleGroupCount)

	log.Debug().
		Str("region", region).
		Float64("firewalls", firewallCount).
		Float64("policies", policyCount).
		Float64("rule_groups", ruleGroupCount).
		Msg("Network Firewall collection completed")

	return nil
}
