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
	mixedInstancesCounts := map[string]float64{"true": 0, "false": 0}
	launchTemplateUsage := make(map[string]float64)
	var lifecycleHookCount float64
	var warmPoolGroupCount float64
	var warmPoolInstanceCount float64
	instanceRefreshCounts := make(map[string]float64)
	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(client, &autoscaling.DescribeAutoScalingGroupsInput{})
	groupNames := make([]string, 0)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, asg := range page.AutoScalingGroups {
			count++
			asgName := aws.ToString(asg.AutoScalingGroupName)
			if asgName == "" {
				continue
			}
			groupNames = append(groupNames, asgName)

			hasMixed := asg.MixedInstancesPolicy != nil
			if hasMixed {
				mixedInstancesCounts["true"]++
			} else {
				mixedInstancesCounts["false"]++
			}

			if asg.LaunchTemplate != nil || asg.MixedInstancesPolicy != nil {
				launchTemplateCount++
			}

			templateID := "none"
			templateVersion := "none"
			if asg.LaunchTemplate != nil {
				templateID = aws.ToString(asg.LaunchTemplate.LaunchTemplateId)
				templateVersion = aws.ToString(asg.LaunchTemplate.Version)
			} else if asg.MixedInstancesPolicy != nil &&
				asg.MixedInstancesPolicy.LaunchTemplate != nil &&
				asg.MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification != nil {
				spec := asg.MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification
				templateID = aws.ToString(spec.LaunchTemplateId)
				templateVersion = aws.ToString(spec.Version)
			}
			if templateID == "" {
				templateID = "unknown"
			}
			if templateVersion == "" {
				templateVersion = "default"
			}
			if templateID != "none" {
				launchTemplateUsage[templateID+"|"+templateVersion]++
			}
		}
	}

	metrics.AutoScalingGroups.WithLabelValues(account, accountName, region).Set(count)
	metrics.AutoScalingGroupsWithLaunchTemplate.WithLabelValues(account, accountName, region).Set(launchTemplateCount)
	for hasMixed, c := range mixedInstancesCounts {
		metrics.AutoScalingGroupsByMixedInstances.WithLabelValues(account, accountName, region, hasMixed).Set(c)
	}
	for key, c := range launchTemplateUsage {
		parts := splitKey(key, 2)
		metrics.AutoScalingLaunchTemplateUsage.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(c)
	}

	var policyCount float64
	policyTypeCounts := make(map[string]float64)
	policyPaginator := autoscaling.NewDescribePoliciesPaginator(client, &autoscaling.DescribePoliciesInput{})
	for policyPaginator.HasMorePages() {
		page, err := policyPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, policy := range page.ScalingPolicies {
			policyCount++
			policyType := aws.ToString(policy.PolicyType)
			if policyType == "" {
				policyType = "SimpleScaling"
			}
			policyTypeCounts[policyType]++
		}
	}
	for policyType, c := range policyTypeCounts {
		metrics.AutoScalingPoliciesByType.WithLabelValues(account, accountName, region, policyType).Set(c)
	}

	for _, groupName := range groupNames {
		hooks, err := client.DescribeLifecycleHooks(ctx, &autoscaling.DescribeLifecycleHooksInput{
			AutoScalingGroupName: aws.String(groupName),
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Str("asg", groupName).Msg("Failed to describe ASG lifecycle hooks")
		} else {
			lifecycleHookCount += float64(len(hooks.LifecycleHooks))
		}

		warmPool, err := client.DescribeWarmPool(ctx, &autoscaling.DescribeWarmPoolInput{
			AutoScalingGroupName: aws.String(groupName),
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Str("asg", groupName).Msg("Failed to describe ASG warm pool")
		} else {
			if len(warmPool.Instances) > 0 || warmPool.WarmPoolConfiguration != nil {
				warmPoolGroupCount++
			}
			warmPoolInstanceCount += float64(len(warmPool.Instances))
		}

		refreshes, err := client.DescribeInstanceRefreshes(ctx, &autoscaling.DescribeInstanceRefreshesInput{
			AutoScalingGroupName: aws.String(groupName),
			MaxRecords:           aws.Int32(10),
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Str("asg", groupName).Msg("Failed to describe ASG instance refreshes")
		} else {
			for _, refresh := range refreshes.InstanceRefreshes {
				status := string(refresh.Status)
				if status == "" {
					status = "unknown"
				}
				instanceRefreshCounts[status]++
			}
		}
	}
	metrics.AutoScalingPolicies.WithLabelValues(account, accountName, region).Set(policyCount)
	metrics.AutoScalingLifecycleHooks.WithLabelValues(account, accountName, region).Set(lifecycleHookCount)
	metrics.AutoScalingWarmPools.WithLabelValues(account, accountName, region).Set(warmPoolGroupCount)
	metrics.AutoScalingWarmPoolInstances.WithLabelValues(account, accountName, region).Set(warmPoolInstanceCount)
	for status, c := range instanceRefreshCounts {
		metrics.AutoScalingInstanceRefreshes.WithLabelValues(account, accountName, region, status).Set(c)
	}

	log.Debug().
		Str("region", region).
		Float64("asg_count", count).
		Float64("asg_with_launch_template_count", launchTemplateCount).
		Float64("policy_count", policyCount).
		Float64("lifecycle_hook_count", lifecycleHookCount).
		Float64("warm_pool_group_count", warmPoolGroupCount).
		Float64("warm_pool_instance_count", warmPoolInstanceCount).
		Int("instance_refresh_status_count", len(instanceRefreshCounts)).
		Msg("Auto Scaling Group collection completed")

	return nil
}
