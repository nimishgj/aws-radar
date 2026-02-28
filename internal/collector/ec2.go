package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type EC2Collector struct{}

func NewEC2Collector() *EC2Collector {
	return &EC2Collector{}
}

func (c *EC2Collector) Name() string {
	return "ec2"
}

func (c *EC2Collector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := ec2.NewFromConfig(cfg)

	// Count instances by type, state, and availability zone
	counts := make(map[string]float64)

	paginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				instanceType := string(instance.InstanceType)
				state := string(instance.State.Name)
				az := aws.ToString(instance.Placement.AvailabilityZone)

				key := instanceType + "|" + state + "|" + az
				counts[key]++
			}
		}
	}

	// Update metrics
	for key, count := range counts {
		parts := splitKey(key, 3)
		metrics.EC2Instances.WithLabelValues(account, accountName, region,
			parts[0], // instance_type
			parts[1], // state
			parts[2], // availability_zone
		).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("instance_combinations", len(counts)).
		Msg("EC2 collection completed")

	return nil
}

// Helper function to split keys
func splitKey(key string, n int) []string {
	result := make([]string, n)
	idx := 0
	start := 0

	for i := 0; i < len(key) && idx < n-1; i++ {
		if key[i] == '|' {
			result[idx] = key[start:i]
			start = i + 1
			idx++
		}
	}
	result[idx] = key[start:]

	return result
}
