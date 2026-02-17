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

func (c *AutoScalingCollector) Collect(ctx context.Context, cfg aws.Config, region, account string) error {
	client := autoscaling.NewFromConfig(cfg)

	var count float64
	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(client, &autoscaling.DescribeAutoScalingGroupsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.AutoScalingGroups))
	}

	metrics.AutoScalingGroups.WithLabelValues(account, region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("asg_count", count).
		Msg("Auto Scaling Group collection completed")

	return nil
}
