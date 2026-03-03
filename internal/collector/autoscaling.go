package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type AutoScalingCollector struct{}

func NewAutoScalingCollector() *AutoScalingCollector {
	return &AutoScalingCollector{}
}

func (c *AutoScalingCollector) Name() string {
	return "autoscaling"
}

func (c *AutoScalingCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := autoscaling.NewFromConfig(cfg)

	var count float64
	var launchTemplateCount float64
	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(client, &autoscaling.DescribeAutoScalingGroupsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, asg := range page.AutoScalingGroups {
			count++
			if asg.LaunchTemplate != nil || asg.MixedInstancesPolicy != nil {
				launchTemplateCount++
			}
		}
	}

	metrics.AutoScalingGroups.WithLabelValues(account, accountName, region).Set(count)
	metrics.AutoScalingGroupsWithLaunchTemplate.WithLabelValues(account, accountName, region).Set(launchTemplateCount)

	var policyCount float64
	policyPaginator := autoscaling.NewDescribePoliciesPaginator(client, &autoscaling.DescribePoliciesInput{})
	for policyPaginator.HasMorePages() {
		page, err := policyPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		policyCount += float64(len(page.ScalingPolicies))
	}
	metrics.AutoScalingPolicies.WithLabelValues(account, accountName, region).Set(policyCount)

	log.Debug().
		Str("region", region).
		Float64("asg_count", count).
		Float64("asg_with_launch_template_count", launchTemplateCount).
		Float64("policy_count", policyCount).
		Msg("Auto Scaling Group collection completed")

	return nil
}
