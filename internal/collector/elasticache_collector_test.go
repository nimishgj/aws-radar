package collector

import (
	"context"
	"testing"

	"github.com/nimishgj/aws-radar/internal/metrics"
)

func TestCollector_elasticache_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewElastiCacheCollector(), "elasticache")
}

func TestCollector_elasticache_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewElastiCacheCollector(), true)
}

func TestCollector_elasticache_Metrics(t *testing.T) {
	server := newMockAWSServer([]mockRoute{
		{
			matcher: queryAction("DescribeCacheClusters"),
			body: `<DescribeCacheClustersResponse>
  <DescribeCacheClustersResult>
    <CacheClusters>
      <CacheCluster>
        <Engine>redis</Engine>
        <CacheNodeType>cache.r6g.large</CacheNodeType>
      </CacheCluster>
      <CacheCluster>
        <Engine>redis</Engine>
        <CacheNodeType>cache.r6g.large</CacheNodeType>
      </CacheCluster>
      <CacheCluster>
        <Engine>memcached</Engine>
        <CacheNodeType>cache.t3.micro</CacheNodeType>
      </CacheCluster>
    </CacheClusters>
  </DescribeCacheClustersResult>
</DescribeCacheClustersResponse>`,
		},
		{
			matcher: queryAction("DescribeServerlessCaches"),
			body: `<DescribeServerlessCachesResponse>
  <DescribeServerlessCachesResult>
    <ServerlessCaches>
      <member>
        <Engine>valkey</Engine>
        <Status>available</Status>
      </member>
    </ServerlessCaches>
  </DescribeServerlessCachesResult>
</DescribeServerlessCachesResponse>`,
		},
		{
			matcher: queryAction("DescribeReplicationGroups"),
			body: `<DescribeReplicationGroupsResponse>
  <DescribeReplicationGroupsResult>
    <ReplicationGroups>
      <ReplicationGroup>
        <Engine>redis</Engine>
        <Status>available</Status>
        <ClusterEnabled>true</ClusterEnabled>
      </ReplicationGroup>
      <ReplicationGroup>
        <Engine>redis</Engine>
        <Status>available</Status>
        <ClusterEnabled>false</ClusterEnabled>
      </ReplicationGroup>
      <ReplicationGroup>
        <Engine>valkey</Engine>
        <Status>creating</Status>
        <ClusterEnabled>true</ClusterEnabled>
      </ReplicationGroup>
    </ReplicationGroups>
  </DescribeReplicationGroupsResult>
</DescribeReplicationGroupsResponse>`,
		},
		{
			matcher: queryAction("DescribeGlobalReplicationGroups"),
			body: `<DescribeGlobalReplicationGroupsResponse>
  <DescribeGlobalReplicationGroupsResult>
    <GlobalReplicationGroups>
      <GlobalReplicationGroup>
        <Status>available</Status>
      </GlobalReplicationGroup>
      <GlobalReplicationGroup>
        <Status>available</Status>
      </GlobalReplicationGroup>
    </GlobalReplicationGroups>
  </DescribeGlobalReplicationGroupsResult>
</DescribeGlobalReplicationGroupsResponse>`,
		},
	})
	defer server.Close()

	metrics.ResetAll()
	cfg := mockAWSConfig(server.URL, "us-east-1")
	err := NewElastiCacheCollector().Collect(context.Background(), cfg, "us-east-1", "123456789012", "test-account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Clusters: 2 redis/cache.r6g.large, 1 memcached/cache.t3.micro
	if v := gaugeValue(metrics.ElastiCacheClusters, "123456789012", "test-account", "us-east-1", "redis", "cache.r6g.large"); v != 2 {
		t.Errorf("ElastiCacheClusters redis/r6g.large: expected 2, got %v", v)
	}
	if v := gaugeValue(metrics.ElastiCacheClusters, "123456789012", "test-account", "us-east-1", "memcached", "cache.t3.micro"); v != 1 {
		t.Errorf("ElastiCacheClusters memcached/t3.micro: expected 1, got %v", v)
	}

	// Serverless caches
	if v := gaugeValue(metrics.ElastiCacheServerlessCaches, "123456789012", "test-account", "us-east-1", "valkey", "available"); v != 1 {
		t.Errorf("ElastiCacheServerlessCaches valkey/available: expected 1, got %v", v)
	}

	// Replication groups
	if v := gaugeValue(metrics.ElastiCacheReplicationGroups, "123456789012", "test-account", "us-east-1", "redis", "available", "true"); v != 1 {
		t.Errorf("ElastiCacheReplicationGroups redis/available/true: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.ElastiCacheReplicationGroups, "123456789012", "test-account", "us-east-1", "redis", "available", "false"); v != 1 {
		t.Errorf("ElastiCacheReplicationGroups redis/available/false: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.ElastiCacheReplicationGroups, "123456789012", "test-account", "us-east-1", "valkey", "creating", "true"); v != 1 {
		t.Errorf("ElastiCacheReplicationGroups valkey/creating/true: expected 1, got %v", v)
	}

	// Global replication groups: 2 available
	if v := gaugeValue(metrics.ElastiCacheGlobalReplicationGroups, "123456789012", "test-account", "us-east-1", "available"); v != 2 {
		t.Errorf("ElastiCacheGlobalReplicationGroups available: expected 2, got %v", v)
	}
}
