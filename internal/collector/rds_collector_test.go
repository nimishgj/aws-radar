package collector

import (
	"context"
	"testing"

	"github.com/nimishgj/aws-radar/internal/metrics"
)

func TestCollector_rds_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewRDSCollector(), "rds")
}

func TestCollector_rds_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewRDSCollector(), true)
}

func TestCollector_rds_Metrics(t *testing.T) {
	server := newMockAWSServer([]mockRoute{
		{
			matcher: queryAction("DescribeDBInstances"),
			body: `<DescribeDBInstancesResponse>
  <DescribeDBInstancesResult>
    <DBInstances>
      <DBInstance>
        <DBInstanceClass>db.t3.micro</DBInstanceClass>
        <Engine>mysql</Engine>
        <EngineVersion>8.0.35</EngineVersion>
        <MultiAZ>false</MultiAZ>
        <DBInstanceStatus>available</DBInstanceStatus>
      </DBInstance>
      <DBInstance>
        <DBInstanceClass>db.r5.large</DBInstanceClass>
        <Engine>postgres</Engine>
        <EngineVersion>15.4</EngineVersion>
        <MultiAZ>true</MultiAZ>
        <DBInstanceStatus>available</DBInstanceStatus>
        <ReadReplicaSourceDBInstanceIdentifier>prod-primary</ReadReplicaSourceDBInstanceIdentifier>
      </DBInstance>
      <DBInstance>
        <DBInstanceClass>db.t3.micro</DBInstanceClass>
        <Engine>mysql</Engine>
        <EngineVersion>8.0.35</EngineVersion>
        <MultiAZ>false</MultiAZ>
        <DBInstanceStatus>available</DBInstanceStatus>
      </DBInstance>
    </DBInstances>
  </DescribeDBInstancesResult>
</DescribeDBInstancesResponse>`,
		},
		{
			matcher: queryAction("DescribeDBProxies"),
			body: `<DescribeDBProxiesResponse>
  <DescribeDBProxiesResult>
    <DBProxies>
      <member>
        <DBProxyName>my-proxy</DBProxyName>
        <EngineFamily>MYSQL</EngineFamily>
        <Status>available</Status>
      </member>
    </DBProxies>
  </DescribeDBProxiesResult>
</DescribeDBProxiesResponse>`,
		},
		{
			matcher: queryAction("DescribeDBProxyTargets"),
			body: `<DescribeDBProxyTargetsResponse>
  <DescribeDBProxyTargetsResult>
    <Targets>
      <member>
        <Type>RDS_INSTANCE</Type>
      </member>
      <member>
        <Type>RDS_INSTANCE</Type>
      </member>
    </Targets>
  </DescribeDBProxyTargetsResult>
</DescribeDBProxyTargetsResponse>`,
		},
		{
			matcher: queryAction("DescribeDBClusters"),
			body: `<DescribeDBClustersResponse>
  <DescribeDBClustersResult>
    <DBClusters>
      <DBCluster>
        <DBClusterIdentifier>aurora-sv2</DBClusterIdentifier>
        <Engine>aurora-mysql</Engine>
        <EngineMode>provisioned</EngineMode>
        <Status>available</Status>
        <ServerlessV2ScalingConfiguration>
          <MinCapacity>0.5</MinCapacity>
          <MaxCapacity>16</MaxCapacity>
        </ServerlessV2ScalingConfiguration>
      </DBCluster>
      <DBCluster>
        <DBClusterIdentifier>aurora-classic</DBClusterIdentifier>
        <Engine>aurora-postgresql</Engine>
        <EngineMode>serverless</EngineMode>
        <Status>available</Status>
      </DBCluster>
    </DBClusters>
  </DescribeDBClustersResult>
</DescribeDBClustersResponse>`,
		},
	})
	defer server.Close()

	metrics.ResetAll()
	cfg := mockAWSConfig(server.URL, "us-east-1")
	err := NewRDSCollector().Collect(context.Background(), cfg, "us-east-1", "123456789012", "test-account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Instance counts: 2x db.t3.micro/mysql/false/available, 1x db.r5.large/postgres/true/available
	if v := gaugeValue(metrics.RDSInstances, "123456789012", "test-account", "us-east-1", "db.t3.micro", "mysql", "false", "available"); v != 2 {
		t.Errorf("RDSInstances mysql/t3.micro: expected 2, got %v", v)
	}
	if v := gaugeValue(metrics.RDSInstances, "123456789012", "test-account", "us-east-1", "db.r5.large", "postgres", "true", "available"); v != 1 {
		t.Errorf("RDSInstances postgres/r5.large: expected 1, got %v", v)
	}

	// Engine version counts
	if v := gaugeValue(metrics.RDSInstancesByEngineVersion, "123456789012", "test-account", "us-east-1", "mysql", "8.0.35"); v != 2 {
		t.Errorf("RDSInstancesByEngineVersion mysql/8.0.35: expected 2, got %v", v)
	}
	if v := gaugeValue(metrics.RDSInstancesByEngineVersion, "123456789012", "test-account", "us-east-1", "postgres", "15.4"); v != 1 {
		t.Errorf("RDSInstancesByEngineVersion postgres/15.4: expected 1, got %v", v)
	}

	// Instance class distribution
	if v := gaugeValue(metrics.RDSInstancesByClass, "123456789012", "test-account", "us-east-1", "db.t3.micro"); v != 2 {
		t.Errorf("RDSInstancesByClass db.t3.micro: expected 2, got %v", v)
	}
	if v := gaugeValue(metrics.RDSInstancesByClass, "123456789012", "test-account", "us-east-1", "db.r5.large"); v != 1 {
		t.Errorf("RDSInstancesByClass db.r5.large: expected 1, got %v", v)
	}

	// Multi-AZ split
	if v := gaugeValue(metrics.RDSInstancesByMultiAZ, "123456789012", "test-account", "us-east-1", "false"); v != 2 {
		t.Errorf("RDSInstancesByMultiAZ false: expected 2, got %v", v)
	}
	if v := gaugeValue(metrics.RDSInstancesByMultiAZ, "123456789012", "test-account", "us-east-1", "true"); v != 1 {
		t.Errorf("RDSInstancesByMultiAZ true: expected 1, got %v", v)
	}

	// Read replicas: 1 postgres replica
	if v := gaugeValue(metrics.RDSReadReplicas, "123456789012", "test-account", "us-east-1", "postgres"); v != 1 {
		t.Errorf("RDSReadReplicas postgres: expected 1, got %v", v)
	}

	// Proxies
	if v := gaugeValue(metrics.RDSProxies, "123456789012", "test-account", "us-east-1", "MYSQL", "available"); v != 1 {
		t.Errorf("RDSProxies MYSQL/available: expected 1, got %v", v)
	}

	// Proxy targets
	if v := gaugeValue(metrics.RDSProxyTargets, "123456789012", "test-account", "us-east-1", "MYSQL", "RDS_INSTANCE"); v != 2 {
		t.Errorf("RDSProxyTargets MYSQL/RDS_INSTANCE: expected 2, got %v", v)
	}

	// Aurora serverless clusters
	if v := gaugeValue(metrics.RDSAuroraServerlessClusters, "123456789012", "test-account", "us-east-1", "aurora-mysql", "provisioned"); v != 1 {
		t.Errorf("RDSAuroraServerlessClusters aurora-mysql/provisioned: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.RDSAuroraServerlessClusters, "123456789012", "test-account", "us-east-1", "aurora-postgresql", "serverless"); v != 1 {
		t.Errorf("RDSAuroraServerlessClusters aurora-postgresql/serverless: expected 1, got %v", v)
	}

	// Aurora serverless by status
	if v := gaugeValue(metrics.RDSAuroraServerlessByStatus, "123456789012", "test-account", "us-east-1", "available"); v != 2 {
		t.Errorf("RDSAuroraServerlessByStatus available: expected 2, got %v", v)
	}

	// Aurora Serverless v2 capacity
	if v := gaugeValue(metrics.RDSAuroraServerlessV2Capacity, "123456789012", "test-account", "us-east-1", "aurora-sv2", "min_capacity_acu"); v != 0.5 {
		t.Errorf("RDSAuroraServerlessV2Capacity min: expected 0.5, got %v", v)
	}
	if v := gaugeValue(metrics.RDSAuroraServerlessV2Capacity, "123456789012", "test-account", "us-east-1", "aurora-sv2", "max_capacity_acu"); v != 16 {
		t.Errorf("RDSAuroraServerlessV2Capacity max: expected 16, got %v", v)
	}
}
