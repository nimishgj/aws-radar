package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type EBSCollector struct{}

func NewEBSCollector() *EBSCollector {
	return &EBSCollector{}
}

func (c *EBSCollector) Name() string {
	return "ebs"
}

func (c *EBSCollector) Collect(ctx context.Context, cfg aws.Config, region, account string) error {
	client := ec2.NewFromConfig(cfg)

	counts := make(map[string]float64)

	paginator := ec2.NewDescribeVolumesPaginator(client, &ec2.DescribeVolumesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, volume := range page.Volumes {
			volumeType := string(volume.VolumeType)
			state := string(volume.State)

			key := volumeType + "|" + state
			counts[key]++
		}
	}

	// Update metrics
	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.EBSVolumes.WithLabelValues(account, region,
			parts[0], // volume_type
			parts[1], // state
		).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("volume_combinations", len(counts)).
		Msg("EBS collection completed")

	return nil
}
