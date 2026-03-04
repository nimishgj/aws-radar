package collector

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/nimishgj/aws-radar/internal/metrics"
)

func TestCollector_ecs_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewECSCollector(), "ecs")
}

func TestCollector_ecs_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewECSCollector(), true)
}

func TestCollector_ecs_Metrics(t *testing.T) {
	server := newMockAWSServer([]mockRoute{
		{
			matcher: jsonTarget("ListClusters"),
			body:    `{"clusterArns":["arn:aws:ecs:us-east-1:123456789012:cluster/prod"]}`,
		},
		{
			matcher: jsonTarget("DescribeClusters"),
			body: `{"clusters":[{
				"clusterArn":"arn:aws:ecs:us-east-1:123456789012:cluster/prod",
				"clusterName":"prod",
				"activeServicesCount":3,
				"runningTasksCount":5,
				"pendingTasksCount":1,
				"registeredContainerInstancesCount":2,
				"defaultCapacityProviderStrategy":[
					{"capacityProvider":"FARGATE","weight":1,"base":1}
				]
			}]}`,
		},
		{
			matcher: jsonTarget("ListServices"),
			body:    `{"serviceArns":["arn:aws:ecs:us-east-1:123456789012:service/prod/web"]}`,
		},
		{
			matcher: jsonTarget("DescribeServices"),
			body: `{"services":[{
				"serviceName":"web",
				"launchType":"FARGATE",
				"status":"ACTIVE"
			}]}`,
		},
		{
			matcher: func(r *http.Request) bool {
				if !strings.Contains(r.Header.Get("X-Amz-Target"), "ListTasks") {
					return false
				}
				bodyBytes, _ := io.ReadAll(r.Body)
				body := string(bodyBytes)
				r.Body = io.NopCloser(strings.NewReader(body))
				if strings.Contains(body, `"RUNNING"`) {
					return true
				}
				return false
			},
			body: `{"taskArns":["arn:aws:ecs:us-east-1:123456789012:task/prod/task1","arn:aws:ecs:us-east-1:123456789012:task/prod/task2"]}`,
		},
		{
			matcher: func(r *http.Request) bool {
				if !strings.Contains(r.Header.Get("X-Amz-Target"), "ListTasks") {
					return false
				}
				bodyBytes, _ := io.ReadAll(r.Body)
				body := string(bodyBytes)
				r.Body = io.NopCloser(strings.NewReader(body))
				return strings.Contains(body, `"PENDING"`)
			},
			body: `{"taskArns":["arn:aws:ecs:us-east-1:123456789012:task/prod/task3"]}`,
		},
		{
			matcher: func(r *http.Request) bool {
				if !strings.Contains(r.Header.Get("X-Amz-Target"), "ListTasks") {
					return false
				}
				return true // STOPPED fallback
			},
			body: `{"taskArns":[]}`,
		},
		{
			matcher: jsonTarget("DescribeTasks"),
			body: `{"tasks":[
				{"taskArn":"task1","launchType":"FARGATE","lastStatus":"RUNNING","desiredStatus":"RUNNING"},
				{"taskArn":"task2","launchType":"EC2","lastStatus":"RUNNING","desiredStatus":"RUNNING"},
				{"taskArn":"task3","launchType":"FARGATE","lastStatus":"PENDING","desiredStatus":"RUNNING"}
			]}`,
		},
		{
			matcher: jsonTarget("DescribeCapacityProviders"),
			body: `{"capacityProviders":[
				{"name":"FARGATE","status":"ACTIVE","type":"FARGATE"},
				{"name":"FARGATE_SPOT","status":"ACTIVE","type":"FARGATE"}
			]}`,
		},
		{
			matcher: func(r *http.Request) bool {
				if !strings.Contains(r.Header.Get("X-Amz-Target"), "ListTaskDefinitions") {
					return false
				}
				bodyBytes, _ := io.ReadAll(r.Body)
				body := string(bodyBytes)
				r.Body = io.NopCloser(strings.NewReader(body))
				return strings.Contains(body, `"ACTIVE"`)
			},
			body: `{"taskDefinitionArns":["arn:aws:ecs:us-east-1:123456789012:task-definition/web:1"]}`,
		},
		{
			matcher: jsonTarget("ListTaskDefinitions"),
			body:    `{"taskDefinitionArns":[]}`,
		},
		{
			matcher: jsonTarget("DescribeTaskDefinition"),
			body: `{"taskDefinition":{
				"taskDefinitionArn":"arn:aws:ecs:us-east-1:123456789012:task-definition/web:1",
				"family":"web",
				"revision":1,
				"runtimePlatform":{"operatingSystemFamily":"LINUX","cpuArchitecture":"ARM64"}
			}}`,
		},
	})
	defer server.Close()

	metrics.ResetAll()
	cfg := mockAWSConfig(server.URL, "us-east-1")
	err := NewECSCollector().Collect(context.Background(), cfg, "us-east-1", "123456789012", "test-account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Cluster depth
	if v := gaugeValue(metrics.ECSClusterDepth, "123456789012", "test-account", "us-east-1", "prod", "active_services"); v != 3 {
		t.Errorf("ECSClusterDepth active_services: expected 3, got %v", v)
	}
	if v := gaugeValue(metrics.ECSClusterDepth, "123456789012", "test-account", "us-east-1", "prod", "running_tasks"); v != 5 {
		t.Errorf("ECSClusterDepth running_tasks: expected 5, got %v", v)
	}
	if v := gaugeValue(metrics.ECSClusterDepth, "123456789012", "test-account", "us-east-1", "prod", "pending_tasks"); v != 1 {
		t.Errorf("ECSClusterDepth pending_tasks: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.ECSClusterDepth, "123456789012", "test-account", "us-east-1", "prod", "container_instances"); v != 2 {
		t.Errorf("ECSClusterDepth container_instances: expected 2, got %v", v)
	}

	// Default capacity provider strategy
	if v := gaugeValue(metrics.ECSDefaultCapacityProviderStrategy, "123456789012", "test-account", "us-east-1", "prod", "FARGATE"); v != 1 {
		t.Errorf("ECSDefaultCapacityProviderStrategy FARGATE: expected 1, got %v", v)
	}

	// Services by launch type
	if v := gaugeValue(metrics.ECSServices, "123456789012", "test-account", "us-east-1", "prod", "FARGATE"); v != 1 {
		t.Errorf("ECSServices prod/FARGATE: expected 1, got %v", v)
	}

	// Services by status
	if v := gaugeValue(metrics.ECSServicesByStatus, "123456789012", "test-account", "us-east-1", "prod", "FARGATE", "ACTIVE"); v != 1 {
		t.Errorf("ECSServicesByStatus ACTIVE: expected 1, got %v", v)
	}

	// Tasks by launch type (FARGATE tasks from RUNNING + PENDING)
	fargateTaskCount := gaugeValue(metrics.ECSTasks, "123456789012", "test-account", "us-east-1", "prod", "FARGATE")
	if fargateTaskCount < 1 {
		t.Errorf("ECSTasks prod/FARGATE: expected >= 1, got %v", fargateTaskCount)
	}

	// Tasks by status: at least one RUNNING/RUNNING FARGATE task
	if v := gaugeValue(metrics.ECSTasksByStatus, "123456789012", "test-account", "us-east-1", "prod", "FARGATE", "RUNNING", "RUNNING"); v < 1 {
		t.Errorf("ECSTasksByStatus FARGATE/RUNNING/RUNNING: expected >= 1, got %v", v)
	}

	// Capacity providers
	if v := gaugeValue(metrics.ECSCapacityProviders, "123456789012", "test-account", "us-east-1", "ACTIVE"); v != 2 {
		t.Errorf("ECSCapacityProviders ACTIVE: expected 2, got %v", v)
	}

	// Capacity providers detailed
	if v := gaugeValue(metrics.ECSCapacityProvidersDetailed, "123456789012", "test-account", "us-east-1", "FARGATE", "ACTIVE"); v != 2 {
		t.Errorf("ECSCapacityProvidersDetailed FARGATE/ACTIVE: expected 2, got %v", v)
	}

	// Task definitions
	if v := gaugeValue(metrics.ECSTaskDefinitions, "123456789012", "test-account", "us-east-1", "ACTIVE"); v != 1 {
		t.Errorf("ECSTaskDefinitions ACTIVE: expected 1, got %v", v)
	}

	// Task definitions detailed
	if v := gaugeValue(metrics.ECSTaskDefinitionsDetailed, "123456789012", "test-account", "us-east-1", "ACTIVE", "web", "1", "LINUX", "ARM64"); v != 1 {
		t.Errorf("ECSTaskDefinitionsDetailed web/1/LINUX/ARM64: expected 1, got %v", v)
	}
}
