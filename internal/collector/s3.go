package collector

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithy "github.com/aws/smithy-go"
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

	regionCounts := make(map[string]float64)
	bpaCounts := make(map[string]map[string]float64) // bucketRegion -> status -> count
	regionClients := map[string]*s3.Client{cfg.Region: client}

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

		// GetPublicAccessBlock requires a regional endpoint, so use a client
		// configured for the bucket's region (cached across buckets in same region).
		bucketClient, ok := regionClients[bucketRegion]
		if !ok {
			regionalCfg := cfg.Copy()
			regionalCfg.Region = bucketRegion
			bucketClient = s3.NewFromConfig(regionalCfg)
			regionClients[bucketRegion] = bucketClient
		}

		// Get Block Public Access configuration
		status := getBucketBPAStatus(ctx, bucketClient, aws.ToString(bucket.Name))
		if bpaCounts[bucketRegion] == nil {
			bpaCounts[bucketRegion] = make(map[string]float64)
		}
		bpaCounts[bucketRegion][status]++
	}

	// Update metrics
	for bucketRegion, count := range regionCounts {
		metrics.S3Buckets.WithLabelValues(account, accountName, bucketRegion).Set(count)
	}
	for bucketRegion, statusCounts := range bpaCounts {
		for status, count := range statusCounts {
			metrics.S3BucketPublicAccessBlock.WithLabelValues(account, accountName, bucketRegion, status).Set(count)
		}
	}

	log.Debug().
		Int("total_buckets", len(output.Buckets)).
		Int("regions", len(regionCounts)).
		Msg("S3 collection completed")

	return nil
}

// getBucketBPAStatus calls GetPublicAccessBlock and classifies the result.
// Returns one of: fully_blocked, partially_blocked, not_blocked, not_set, error.
func getBucketBPAStatus(ctx context.Context, client *s3.Client, bucket string) string {
	out, err := client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NoSuchPublicAccessBlockConfiguration" {
			return "not_set"
		}
		log.Warn().
			Err(err).
			Str("bucket", bucket).
			Msg("Failed to get bucket public access block")
		return "error"
	}
	return classifyBPA(out.PublicAccessBlockConfiguration)
}

func classifyBPA(c *s3types.PublicAccessBlockConfiguration) string {
	if c == nil {
		return "not_set"
	}
	flags := []bool{
		aws.ToBool(c.BlockPublicAcls),
		aws.ToBool(c.IgnorePublicAcls),
		aws.ToBool(c.BlockPublicPolicy),
		aws.ToBool(c.RestrictPublicBuckets),
	}
	trueCount := 0
	for _, f := range flags {
		if f {
			trueCount++
		}
	}
	switch trueCount {
	case 4:
		return "fully_blocked"
	case 0:
		return "not_blocked"
	default:
		return "partially_blocked"
	}
}
