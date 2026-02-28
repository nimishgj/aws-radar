package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type ElastiCacheCollector struct{}

func NewElastiCacheCollector() *ElastiCacheCollector {
	return &ElastiCacheCollector{}
}

func (c *ElastiCacheCollector) Name() string {
	return "elasticache"
}

func (c *ElastiCacheCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := elasticache.NewFromConfig(cfg)

	counts := make(map[string]float64)

	paginator := elasticache.NewDescribeCacheClustersPaginator(client, &elasticache.DescribeCacheClustersInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, cluster := range page.CacheClusters {
			engine := aws.ToString(cluster.Engine)
			nodeType := aws.ToString(cluster.CacheNodeType)

			key := engine + "|" + nodeType
			counts[key]++
		}
	}

	// Update metrics
	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.ElastiCacheClusters.WithLabelValues(account, accountName, region,
			parts[0], // engine
			parts[1], // cache_node_type
		).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("cluster_combinations", len(counts)).
		Msg("ElastiCache collection completed")

	return nil
}
