package collector

import (
	"context"
	"strconv"

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
	serverlessCounts := make(map[string]float64)
	replicationGroupCounts := make(map[string]float64)
	globalReplicationGroupCounts := make(map[string]float64)

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

	serverlessPaginator := elasticache.NewDescribeServerlessCachesPaginator(client, &elasticache.DescribeServerlessCachesInput{})
	for serverlessPaginator.HasMorePages() {
		page, err := serverlessPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, cache := range page.ServerlessCaches {
			engine := aws.ToString(cache.Engine)
			if engine == "" {
				engine = "unknown"
			}
			status := aws.ToString(cache.Status)
			if status == "" {
				status = "unknown"
			}
			key := engine + "|" + status
			serverlessCounts[key]++
		}
	}

	replicationGroupPaginator := elasticache.NewDescribeReplicationGroupsPaginator(client, &elasticache.DescribeReplicationGroupsInput{})
	for replicationGroupPaginator.HasMorePages() {
		page, err := replicationGroupPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, group := range page.ReplicationGroups {
			engine := aws.ToString(group.Engine)
			if engine == "" {
				engine = "unknown"
			}
			status := aws.ToString(group.Status)
			if status == "" {
				status = "unknown"
			}
			clusterEnabled := "false"
			if group.ClusterEnabled != nil {
				clusterEnabled = strconv.FormatBool(*group.ClusterEnabled)
			}
			key := engine + "|" + status + "|" + clusterEnabled
			replicationGroupCounts[key]++
		}
	}

	globalReplicationGroupPaginator := elasticache.NewDescribeGlobalReplicationGroupsPaginator(client, &elasticache.DescribeGlobalReplicationGroupsInput{})
	for globalReplicationGroupPaginator.HasMorePages() {
		page, err := globalReplicationGroupPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, group := range page.GlobalReplicationGroups {
			status := aws.ToString(group.Status)
			if status == "" {
				status = "unknown"
			}
			globalReplicationGroupCounts[status]++
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

	for key, count := range serverlessCounts {
		parts := splitKey(key, 2)
		metrics.ElastiCacheServerlessCaches.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	for key, count := range replicationGroupCounts {
		parts := splitKey(key, 3)
		metrics.ElastiCacheReplicationGroups.WithLabelValues(account, accountName, region, parts[0], parts[1], parts[2]).Set(count)
	}

	for status, count := range globalReplicationGroupCounts {
		metrics.ElastiCacheGlobalReplicationGroups.WithLabelValues(account, accountName, region, status).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("cluster_combinations", len(counts)).
		Int("serverless_combinations", len(serverlessCounts)).
		Int("replication_group_combinations", len(replicationGroupCounts)).
		Int("global_replication_group_statuses", len(globalReplicationGroupCounts)).
		Msg("ElastiCache collection completed")

	return nil
}
