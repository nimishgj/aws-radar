package collector

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type RDSCollector struct{}

func NewRDSCollector() *RDSCollector {
	return &RDSCollector{}
}

func (c *RDSCollector) Name() string {
	return "rds"
}

func (c *RDSCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := rds.NewFromConfig(cfg)

	counts := make(map[string]float64)

	paginator := rds.NewDescribeDBInstancesPaginator(client, &rds.DescribeDBInstancesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, instance := range page.DBInstances {
			dbClass := aws.ToString(instance.DBInstanceClass)
			engine := aws.ToString(instance.Engine)
			multiAZ := strconv.FormatBool(instance.MultiAZ != nil && *instance.MultiAZ)
			status := aws.ToString(instance.DBInstanceStatus)
			if status == "" {
				status = "unknown"
			}

			key := dbClass + "|" + engine + "|" + multiAZ + "|" + status
			counts[key]++
		}
	}

	proxyCounts := make(map[string]float64)
	proxyPaginator := rds.NewDescribeDBProxiesPaginator(client, &rds.DescribeDBProxiesInput{})
	for proxyPaginator.HasMorePages() {
		page, err := proxyPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, proxy := range page.DBProxies {
			engineFamily := aws.ToString(proxy.EngineFamily)
			if engineFamily == "" {
				engineFamily = "unknown"
			}
			status := string(proxy.Status)
			if status == "" {
				status = "unknown"
			}
			key := engineFamily + "|" + status
			proxyCounts[key]++
		}
	}

	auroraServerlessCounts := make(map[string]float64)
	clusterPaginator := rds.NewDescribeDBClustersPaginator(client, &rds.DescribeDBClustersInput{})
	for clusterPaginator.HasMorePages() {
		page, err := clusterPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, cluster := range page.DBClusters {
			serverless := false
			engineMode := aws.ToString(cluster.EngineMode)
			if engineMode == "serverless" {
				serverless = true
			}
			if cluster.ServerlessV2ScalingConfiguration != nil {
				serverless = true
			}
			if !serverless {
				continue
			}
			engine := aws.ToString(cluster.Engine)
			if engine == "" {
				engine = "unknown"
			}
			if engineMode == "" {
				engineMode = "provisioned"
			}
			key := engine + "|" + engineMode
			auroraServerlessCounts[key]++
		}
	}

	// Update metrics
	for key, count := range counts {
		parts := splitKey(key, 4)
		metrics.RDSInstances.WithLabelValues(account, accountName, region,
			parts[0], // db_instance_class
			parts[1], // engine
			parts[2], // multi_az
			parts[3], // status
		).Set(count)
	}

	for key, count := range proxyCounts {
		parts := splitKey(key, 2)
		metrics.RDSProxies.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	for key, count := range auroraServerlessCounts {
		parts := splitKey(key, 2)
		metrics.RDSAuroraServerlessClusters.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("instance_combinations", len(counts)).
		Int("proxy_combinations", len(proxyCounts)).
		Int("aurora_serverless_combinations", len(auroraServerlessCounts)).
		Msg("RDS collection completed")

	return nil
}
