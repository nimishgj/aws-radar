package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
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

	log.Debug().
		Str("region", region).
		Int("service_combinations", len(serviceCounts)).
		Int("task_combinations", len(taskCounts)).
		Msg("ECS collection completed")

	return nil
}
