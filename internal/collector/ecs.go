package collector

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type ECSCollector struct{}

func NewECSCollector() *ECSCollector {
	return &ECSCollector{}
}

func (c *ECSCollector) Name() string {
	return "ecs"
}

func (c *ECSCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := ecs.NewFromConfig(cfg)

	// List clusters
	clusterPaginator := ecs.NewListClustersPaginator(client, &ecs.ListClustersInput{})

	serviceCounts := make(map[string]float64)
	taskCounts := make(map[string]float64)
	serviceStatusCounts := make(map[string]float64)
	capacityProviderCounts := make(map[string]float64)
	capacityProviderDetailedCounts := make(map[string]float64)
	defaultCapacityProviderStrategyCounts := make(map[string]float64)
	taskDefinitionCounts := make(map[string]float64)
	taskDefinitionDetailedCounts := make(map[string]float64)
	taskStatusCounts := make(map[string]float64)
	clusterDepth := make(map[string]float64)

	for clusterPaginator.HasMorePages() {
		clusterPage, err := clusterPaginator.NextPage(ctx)
		if err != nil {
			return err
		}

		if len(clusterPage.ClusterArns) == 0 {
			continue
		}

		// Describe clusters to get names
		describeOutput, err := client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
			Clusters: clusterPage.ClusterArns,
		})
		if err != nil {
			return err
		}

		for _, cluster := range describeOutput.Clusters {
			clusterName := aws.ToString(cluster.ClusterName)
			if clusterName == "" {
				clusterName = "unknown"
			}

			clusterDepth[clusterName+"|active_services"] = float64(cluster.ActiveServicesCount)
			clusterDepth[clusterName+"|running_tasks"] = float64(cluster.RunningTasksCount)
			clusterDepth[clusterName+"|pending_tasks"] = float64(cluster.PendingTasksCount)
			clusterDepth[clusterName+"|container_instances"] = float64(cluster.RegisteredContainerInstancesCount)

			for _, strategy := range cluster.DefaultCapacityProviderStrategy {
				provider := aws.ToString(strategy.CapacityProvider)
				if provider == "" {
					provider = "unknown"
				}
				key := clusterName + "|" + provider
				defaultCapacityProviderStrategyCounts[key]++
			}

			// Count services in this cluster
			servicePaginator := ecs.NewListServicesPaginator(client, &ecs.ListServicesInput{
				Cluster: cluster.ClusterArn,
			})

			for servicePaginator.HasMorePages() {
				servicePage, err := servicePaginator.NextPage(ctx)
				if err != nil {
					log.Warn().
						Err(err).
						Str("cluster", clusterName).
						Msg("Failed to list services")
					break
				}

				if len(servicePage.ServiceArns) > 0 {
					// Describe services to get launch type
					descServicesOutput, err := client.DescribeServices(ctx, &ecs.DescribeServicesInput{
						Cluster:  cluster.ClusterArn,
						Services: servicePage.ServiceArns,
					})
					if err != nil {
						log.Warn().Err(err).Msg("Failed to describe services")
						continue
					}

					for _, service := range descServicesOutput.Services {
						launchType := string(service.LaunchType)
						if launchType == "" {
							launchType = "EC2"
						}
						key := clusterName + "|" + launchType
						serviceCounts[key]++

						status := aws.ToString(service.Status)
						if status == "" {
							status = "UNKNOWN"
						}
						statusKey := clusterName + "|" + launchType + "|" + status
						serviceStatusCounts[statusKey]++
					}
				}
			}

			for _, desiredStatus := range []string{"RUNNING", "PENDING", "STOPPED"} {
				taskPaginator := ecs.NewListTasksPaginator(client, &ecs.ListTasksInput{
					Cluster:       cluster.ClusterArn,
					DesiredStatus: ecsTypes.DesiredStatus(desiredStatus),
				})

				for taskPaginator.HasMorePages() {
					taskPage, err := taskPaginator.NextPage(ctx)
					if err != nil {
						log.Warn().
							Err(err).
							Str("cluster", clusterName).
							Str("desired_status", desiredStatus).
							Msg("Failed to list tasks")
						break
					}

					if len(taskPage.TaskArns) > 0 {
						descTasksOutput, err := client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
							Cluster: cluster.ClusterArn,
							Tasks:   taskPage.TaskArns,
						})
						if err != nil {
							log.Warn().Err(err).Msg("Failed to describe tasks")
							continue
						}

						for _, task := range descTasksOutput.Tasks {
							launchType := string(task.LaunchType)
							if launchType == "" {
								launchType = "EC2"
							}
							key := clusterName + "|" + launchType
							taskCounts[key]++

							lastStatus := aws.ToString(task.LastStatus)
							if lastStatus == "" {
								lastStatus = "UNKNOWN"
							}
							desired := aws.ToString(task.DesiredStatus)
							if desired == "" {
								desired = desiredStatus
							}
							statusKey := clusterName + "|" + launchType + "|" + lastStatus + "|" + desired
							taskStatusCounts[statusKey]++
						}
					}
				}
			}
		}
	}

	// Capacity providers (account-level in region)
	var capNextToken *string
	for {
		page, err := client.DescribeCapacityProviders(ctx, &ecs.DescribeCapacityProvidersInput{
			NextToken:  capNextToken,
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to describe ECS capacity providers")
			break
		}
		for _, cp := range page.CapacityProviders {
			status := string(cp.Status)
			if status == "" {
				status = "UNKNOWN"
			}
			capacityProviderCounts[status]++

			cpType := string(cp.Type)
			if cpType == "" {
				cpType = "UNKNOWN"
			}
			detailedKey := cpType + "|" + status
			capacityProviderDetailedCounts[detailedKey]++
		}
		if page.NextToken == nil || *page.NextToken == "" {
			break
		}
		capNextToken = page.NextToken
	}

	// Task definitions by status
	for _, tdStatus := range []ecsTypes.TaskDefinitionStatus{ecsTypes.TaskDefinitionStatusActive, ecsTypes.TaskDefinitionStatusInactive, ecsTypes.TaskDefinitionStatusDeleteInProgress} {
		tdPaginator := ecs.NewListTaskDefinitionsPaginator(client, &ecs.ListTaskDefinitionsInput{
			Status: tdStatus,
		})
		for tdPaginator.HasMorePages() {
			page, err := tdPaginator.NextPage(ctx)
			if err != nil {
				log.Warn().Err(err).Str("region", region).Str("status", string(tdStatus)).Msg("Failed to list ECS task definitions")
				break
			}
			taskDefinitionCounts[string(tdStatus)] += float64(len(page.TaskDefinitionArns))

			for _, tdArn := range page.TaskDefinitionArns {
				desc, descErr := client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
					TaskDefinition: aws.String(tdArn),
				})
				if descErr != nil {
					log.Warn().Err(descErr).Str("region", region).Str("task_definition", tdArn).Msg("Failed to describe ECS task definition")
					continue
				}

				family := aws.ToString(desc.TaskDefinition.Family)
				if family == "" {
					family = "unknown"
				}
				revision := "0"
				revision = strconv.FormatInt(int64(desc.TaskDefinition.Revision), 10)
				osFamily := "unknown"
				cpuArch := "unknown"
				if desc.TaskDefinition.RuntimePlatform != nil {
					osFamily = string(desc.TaskDefinition.RuntimePlatform.OperatingSystemFamily)
					cpuArch = string(desc.TaskDefinition.RuntimePlatform.CpuArchitecture)
					if osFamily == "" {
						osFamily = "unknown"
					}
					if cpuArch == "" {
						cpuArch = "unknown"
					}
				}
				detailedKey := string(tdStatus) + "|" + family + "|" + revision + "|" + osFamily + "|" + cpuArch
				taskDefinitionDetailedCounts[detailedKey]++
			}
		}
	}

	// Update service metrics
	for key, count := range serviceCounts {
		parts := splitKey(key, 2)
		metrics.ECSServices.WithLabelValues(account, accountName, region,
			parts[0], // cluster_name
			parts[1], // launch_type
		).Set(count)
	}

	// Update task metrics
	for key, count := range taskCounts {
		parts := splitKey(key, 2)
		metrics.ECSTasks.WithLabelValues(account, accountName, region,
			parts[0], // cluster_name
			parts[1], // launch_type
		).Set(count)
	}

	for key, count := range serviceStatusCounts {
		parts := splitKey(key, 3)
		metrics.ECSServicesByStatus.WithLabelValues(account, accountName, region,
			parts[0], // cluster_name
			parts[1], // launch_type
			parts[2], // status
		).Set(count)
	}

	for key, count := range taskStatusCounts {
		parts := splitKey(key, 4)
		metrics.ECSTasksByStatus.WithLabelValues(account, accountName, region,
			parts[0], // cluster_name
			parts[1], // launch_type
			parts[2], // last_status
			parts[3], // desired_status
		).Set(count)
	}

	for key, count := range clusterDepth {
		parts := splitKey(key, 2)
		metrics.ECSClusterDepth.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	for status, count := range capacityProviderCounts {
		metrics.ECSCapacityProviders.WithLabelValues(account, accountName, region, status).Set(count)
	}

	for key, count := range capacityProviderDetailedCounts {
		parts := splitKey(key, 2)
		metrics.ECSCapacityProvidersDetailed.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	for key, count := range defaultCapacityProviderStrategyCounts {
		parts := splitKey(key, 2)
		metrics.ECSDefaultCapacityProviderStrategy.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}

	for status, count := range taskDefinitionCounts {
		metrics.ECSTaskDefinitions.WithLabelValues(account, accountName, region, status).Set(count)
	}

	for key, count := range taskDefinitionDetailedCounts {
		parts := splitKey(key, 5)
		metrics.ECSTaskDefinitionsDetailed.WithLabelValues(account, accountName, region,
			parts[0], // status
			parts[1], // family
			parts[2], // revision
			parts[3], // os_family
			parts[4], // cpu_architecture
		).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("service_combinations", len(serviceCounts)).
		Int("task_combinations", len(taskCounts)).
		Int("service_status_combinations", len(serviceStatusCounts)).
		Int("task_status_combinations", len(taskStatusCounts)).
		Int("capacity_provider_statuses", len(capacityProviderCounts)).
		Int("capacity_provider_detailed_combinations", len(capacityProviderDetailedCounts)).
		Int("default_capacity_provider_strategy_combinations", len(defaultCapacityProviderStrategyCounts)).
		Int("task_definition_statuses", len(taskDefinitionCounts)).
		Int("task_definition_detailed_combinations", len(taskDefinitionDetailedCounts)).
		Msg("ECS collection completed")

	return nil
}
