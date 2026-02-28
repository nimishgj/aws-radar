package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type S3Collector struct{}

func NewS3Collector() *S3Collector {
	return &S3Collector{}
}

func (c *S3Collector) Name() string {
	return "s3"
}

func (c *S3Collector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := s3.NewFromConfig(cfg)

	// S3 ListBuckets returns all buckets regardless of region
	// We only count in one region to avoid duplicates
	if region != cfg.Region {
		// Only collect in the primary region
		return nil
	}

	output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	// Count buckets by region
	regionCounts := make(map[string]float64)

	for _, bucket := range output.Buckets {
		// Get bucket location
		locOutput, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			log.Warn().
				Err(err).
				Str("bucket", aws.ToString(bucket.Name)).
				Msg("Failed to get bucket location")
			continue
		}

		bucketRegion := string(locOutput.LocationConstraint)
		if bucketRegion == "" {
			bucketRegion = "us-east-1" // Default region
		}

		regionCounts[bucketRegion]++
	}

	// Update metrics
	for bucketRegion, count := range regionCounts {
		metrics.S3Buckets.WithLabelValues(account, accountName, bucketRegion).Set(count)
	}

	log.Debug().
		Int("total_buckets", len(output.Buckets)).
		Int("regions", len(regionCounts)).
		Msg("S3 collection completed")

	return nil
}
