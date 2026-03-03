package collector

import (
	"context"

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
	taskDefinitionCounts := make(map[string]float64)

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

			// Count tasks in this cluster
			taskPaginator := ecs.NewListTasksPaginator(client, &ecs.ListTasksInput{
				Cluster: cluster.ClusterArn,
			})

			for taskPaginator.HasMorePages() {
				taskPage, err := taskPaginator.NextPage(ctx)
				if err != nil {
					log.Warn().
						Err(err).
						Str("cluster", clusterName).
						Msg("Failed to list tasks")
					break
				}

				if len(taskPage.TaskArns) > 0 {
					// Describe tasks to get launch type
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

	for status, count := range capacityProviderCounts {
		metrics.ECSCapacityProviders.WithLabelValues(account, accountName, region, status).Set(count)
	}

	for status, count := range taskDefinitionCounts {
		metrics.ECSTaskDefinitions.WithLabelValues(account, accountName, region, status).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("service_combinations", len(serviceCounts)).
		Int("task_combinations", len(taskCounts)).
		Int("service_status_combinations", len(serviceStatusCounts)).
		Int("capacity_provider_statuses", len(capacityProviderCounts)).
		Int("task_definition_statuses", len(taskDefinitionCounts)).
		Msg("ECS collection completed")

	return nil
}
