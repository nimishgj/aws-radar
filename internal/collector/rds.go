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
	engineVersionCounts := make(map[string]float64)
	instanceClassCounts := make(map[string]float64)
	multiAZCounts := make(map[string]float64)
	readReplicaCounts := make(map[string]float64)

	paginator := rds.NewDescribeDBInstancesPaginator(client, &rds.DescribeDBInstancesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, instance := range page.DBInstances {
			dbClass := aws.ToString(instance.DBInstanceClass)
			engine := aws.ToString(instance.Engine)
			engineVersion := aws.ToString(instance.EngineVersion)
			if engineVersion == "" {
				engineVersion = "unknown"
			}
			if engine == "" {
				engine = "unknown"
			}
			multiAZ := strconv.FormatBool(instance.MultiAZ != nil && *instance.MultiAZ)
			status := aws.ToString(instance.DBInstanceStatus)
			if status == "" {
				status = "unknown"
			}

			key := dbClass + "|" + engine + "|" + multiAZ + "|" + status
			counts[key]++
			engineVersionCounts[engine+"|"+engineVersion]++
			instanceClassCounts[dbClass]++
			multiAZCounts[multiAZ]++
			if aws.ToString(instance.ReadReplicaSourceDBInstanceIdentifier) != "" {
				readReplicaCounts[engine]++
			}
		}
	}

	proxyCounts := make(map[string]float64)
	proxyTargetCounts := make(map[string]float64)
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

			targetPaginator := rds.NewDescribeDBProxyTargetsPaginator(client, &rds.DescribeDBProxyTargetsInput{
				DBProxyName: proxy.DBProxyName,
			})
			for targetPaginator.HasMorePages() {
				targetPage, targetErr := targetPaginator.NextPage(ctx)
				if targetErr != nil {
					log.Warn().Err(targetErr).Str("region", region).Str("proxy", aws.ToString(proxy.DBProxyName)).Msg("Failed to describe RDS proxy targets")
					break
				}
				for _, target := range targetPage.Targets {
					targetType := string(target.Type)
					if targetType == "" {
						targetType = "unknown"
					}
					targetKey := engineFamily + "|" + targetType
					proxyTargetCounts[targetKey]++
				}
			}
		}
	}

	auroraServerlessCounts := make(map[string]float64)
	auroraServerlessStatusCounts := make(map[string]float64)
	auroraServerlessV2Capacity := make(map[string]map[string]float64)
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
			status := aws.ToString(cluster.Status)
			if status == "" {
				status = "unknown"
			}
			if engineMode == "" {
				engineMode = "provisioned"
			}
			key := engine + "|" + engineMode
			auroraServerlessCounts[key]++
			auroraServerlessStatusCounts[status]++

			if cluster.ServerlessV2ScalingConfiguration != nil {
				clusterID := aws.ToString(cluster.DBClusterIdentifier)
				if clusterID == "" {
					clusterID = "unknown"
				}
				if _, ok := auroraServerlessV2Capacity[clusterID]; !ok {
					auroraServerlessV2Capacity[clusterID] = map[string]float64{}
				}
				if cluster.ServerlessV2ScalingConfiguration.MinCapacity != nil {
					auroraServerlessV2Capacity[clusterID]["min_capacity_acu"] = *cluster.ServerlessV2ScalingConfiguration.MinCapacity
				}
				if cluster.ServerlessV2ScalingConfiguration.MaxCapacity != nil {
					auroraServerlessV2Capacity[clusterID]["max_capacity_acu"] = *cluster.ServerlessV2ScalingConfiguration.MaxCapacity
				}
			}
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

	for key, count := range proxyTargetCounts {
		parts := splitKey(key, 2)
		metrics.RDSProxyTargets.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	for key, count := range auroraServerlessCounts {
		parts := splitKey(key, 2)
		metrics.RDSAuroraServerlessClusters.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	for key, count := range engineVersionCounts {
		parts := splitKey(key, 2)
		metrics.RDSInstancesByEngineVersion.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	for dbClass, count := range instanceClassCounts {
		metrics.RDSInstancesByClass.WithLabelValues(account, accountName, region, dbClass).Set(count)
	}

	for multiAZ, count := range multiAZCounts {
		metrics.RDSInstancesByMultiAZ.WithLabelValues(account, accountName, region, multiAZ).Set(count)
	}

	for engine, count := range readReplicaCounts {
		metrics.RDSReadReplicas.WithLabelValues(account, accountName, region, engine).Set(count)
	}

	for status, count := range auroraServerlessStatusCounts {
		metrics.RDSAuroraServerlessByStatus.WithLabelValues(account, accountName, region, status).Set(count)
	}

	for clusterID, values := range auroraServerlessV2Capacity {
		for metricName, value := range values {
			metrics.RDSAuroraServerlessV2Capacity.WithLabelValues(account, accountName, region, clusterID, metricName).Set(value)
		}
	}

	log.Debug().
		Str("region", region).
		Int("instance_combinations", len(counts)).
		Int("proxy_combinations", len(proxyCounts)).
		Int("proxy_target_combinations", len(proxyTargetCounts)).
		Int("aurora_serverless_combinations", len(auroraServerlessCounts)).
		Int("aurora_serverless_statuses", len(auroraServerlessStatusCounts)).
		Msg("RDS collection completed")

	return nil
}
