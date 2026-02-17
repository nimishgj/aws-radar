package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type EKSCollector struct{}

func NewEKSCollector() *EKSCollector {
	return &EKSCollector{}
}

func (c *EKSCollector) Name() string {
	return "eks"
}

func (c *EKSCollector) Collect(ctx context.Context, cfg aws.Config, region, account string) error {
	client := eks.NewFromConfig(cfg)

	counts := make(map[string]float64)

	paginator := eks.NewListClustersPaginator(client, &eks.ListClustersInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, clusterName := range page.Clusters {
			// Get cluster details
			descOutput, err := client.DescribeCluster(ctx, &eks.DescribeClusterInput{
				Name: aws.String(clusterName),
			})
			if err != nil {
				log.Warn().
					Err(err).
					Str("cluster", clusterName).
					Msg("Failed to describe EKS cluster")
				continue
			}

			version := aws.ToString(descOutput.Cluster.Version)
			if version == "" {
				version = "unknown"
			}
			status := string(descOutput.Cluster.Status)
			if status == "" {
				status = "unknown"
			}

			key := version + "|" + status
			counts[key]++
		}
	}

	// Update metrics
	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.EKSClusters.WithLabelValues(account, region, parts[0], parts[1]).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("cluster_versions", len(counts)).
		Msg("EKS collection completed")

	return nil
}
