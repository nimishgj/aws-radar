package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// API Gateway Metrics
	APIGatewayRestAPIs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_apigateway_rest_apis_total",
			Help: "Total number of API Gateway REST APIs",
		},
		[]string{"region"},
	)

	APIGatewayV2APIs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_apigatewayv2_apis_total",
			Help: "Total number of API Gateway v2 APIs",
		},
		[]string{"region", "protocol"},
	)

	// Auto Scaling Metrics
	AutoScalingGroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_groups_total",
			Help: "Total number of Auto Scaling Groups",
		},
		[]string{"region"},
	)

	// Athena Metrics
	AthenaWorkgroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_athena_workgroups_total",
			Help: "Total number of Athena workgroups",
		},
		[]string{"region"},
	)

	// ECR Metrics
	ECRRepositories = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecr_repositories_total",
			Help: "Total number of ECR repositories",
		},
		[]string{"region"},
	)

	// EC2 Metrics
	EC2Instances = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ec2_instances_total",
			Help: "Total number of EC2 instances",
		},
		[]string{"region", "instance_type", "state", "availability_zone"},
	)

	// S3 Metrics
	S3Buckets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_s3_buckets_total",
			Help: "Total number of S3 buckets",
		},
		[]string{"region"},
	)

	// RDS Metrics
	RDSInstances = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_instances_total",
			Help: "Total number of RDS instances",
		},
		[]string{"region", "db_instance_class", "engine", "multi_az", "status"},
	)

	// Lambda Metrics
	LambdaFunctions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_lambda_functions_total",
			Help: "Total number of Lambda functions",
		},
		[]string{"region", "runtime", "memory_size"},
	)

	// ECS Metrics
	ECSServices = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_services_total",
			Help: "Total number of ECS services",
		},
		[]string{"region", "cluster_name", "launch_type"},
	)

	ECSTasks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_tasks_total",
			Help: "Total number of ECS tasks",
		},
		[]string{"region", "cluster_name", "launch_type"},
	)

	// EKS Metrics
	EKSClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_eks_clusters_total",
			Help: "Total number of EKS clusters",
		},
		[]string{"region", "version", "status"},
	)

	// ELB Metrics
	ELBClassic = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elb_classic_total",
			Help: "Total number of Classic Load Balancers",
		},
		[]string{"region", "scheme"},
	)

	ELBV2 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elbv2_total",
			Help: "Total number of ALB/NLB load balancers",
		},
		[]string{"region", "type", "scheme"},
	)

	// DynamoDB Metrics
	DynamoDBTables = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_dynamodb_tables_total",
			Help: "Total number of DynamoDB tables",
		},
		[]string{"region", "billing_mode"},
	)

	// ElastiCache Metrics
	ElastiCacheClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elasticache_clusters_total",
			Help: "Total number of ElastiCache clusters",
		},
		[]string{"region", "engine", "cache_node_type"},
	)

	// SQS Metrics
	SQSQueues = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sqs_queues_total",
			Help: "Total number of SQS queues",
		},
		[]string{"region", "queue_type"},
	)

	// SNS Metrics
	SNSTopics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sns_topics_total",
			Help: "Total number of SNS topics",
		},
		[]string{"region"},
	)

	// CloudFront Metrics
	CloudFrontDistributions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cloudfront_distributions_total",
			Help: "Total number of CloudFront distributions",
		},
		[]string{"price_class", "enabled"},
	)

	// EBS Metrics
	EBSVolumes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ebs_volumes_total",
			Help: "Total number of EBS volumes",
		},
		[]string{"region", "volume_type", "state"},
	)

	// VPC Metrics
	VPCs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_vpc_total",
			Help: "Total number of VPCs",
		},
		[]string{"region", "state"},
	)

	Subnets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_subnet_total",
			Help: "Total number of subnets",
		},
		[]string{"region", "availability_zone"},
	)

	SecurityGroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_security_groups_total",
			Help: "Total number of security groups",
		},
		[]string{"region", "vpc_id"},
	)

	NATGateways = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_nat_gateways_total",
			Help: "Total number of NAT gateways",
		},
		[]string{"region", "state"},
	)

	InternetGateways = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_internet_gateways_total",
			Help: "Total number of Internet gateways",
		},
		[]string{"region"},
	)

	// EFS Metrics
	EFSFileSystems = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_efs_filesystems_total",
			Help: "Total number of EFS file systems",
		},
		[]string{"region"},
	)

	// EventBridge Metrics
	EventBridgeRules = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_eventbridge_rules_total",
			Help: "Total number of EventBridge rules",
		},
		[]string{"region", "event_bus"},
	)

	// Glue Metrics
	GlueJobs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_glue_jobs_total",
			Help: "Total number of Glue jobs",
		},
		[]string{"region"},
	)

	// OpenSearch Metrics
	OpenSearchDomains = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_opensearch_domains_total",
			Help: "Total number of OpenSearch domains",
		},
		[]string{"region"},
	)

	// Secrets Manager Metrics
	SecretsManagerSecrets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_secretsmanager_secrets_total",
			Help: "Total number of Secrets Manager secrets",
		},
		[]string{"region"},
	)

	// SSM Metrics
	SSMParameters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ssm_parameters_total",
			Help: "Total number of SSM parameters",
		},
		[]string{"region", "type"},
	)

	// Step Functions Metrics
	SFNStateMachines = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sfn_state_machines_total",
			Help: "Total number of Step Functions state machines",
		},
		[]string{"region", "type"},
	)

	// Route53 Metrics
	Route53HostedZones = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_route53_hosted_zones_total",
			Help: "Total number of Route53 hosted zones",
		},
		[]string{},
	)

	// ACM Metrics
	ACMCertificates = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_acm_certificates_total",
			Help: "Total number of ACM certificates",
		},
		[]string{"region", "status", "type"},
	)

	// IAM Metrics
	IAMUsers = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_iam_users_total",
			Help: "Total number of IAM users",
		},
		[]string{},
	)

	IAMRoles = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_iam_roles_total",
			Help: "Total number of IAM roles",
		},
		[]string{},
	)

	// Collection Metrics
	CollectionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "aws_radar_collection_duration_seconds",
			Help:    "Duration of AWS resource collection",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"collector"},
	)

	CollectionErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aws_radar_collection_errors_total",
			Help: "Total number of collection errors",
		},
		[]string{"collector", "region"},
	)
)

// ResetAll resets all gauge metrics before a new collection cycle
func ResetAll() {
	APIGatewayRestAPIs.Reset()
	APIGatewayV2APIs.Reset()
	AutoScalingGroups.Reset()
	AthenaWorkgroups.Reset()
	ECRRepositories.Reset()
	EC2Instances.Reset()
	EFSFileSystems.Reset()
	EventBridgeRules.Reset()
	GlueJobs.Reset()
	S3Buckets.Reset()
	RDSInstances.Reset()
	LambdaFunctions.Reset()
	ECSServices.Reset()
	ECSTasks.Reset()
	EKSClusters.Reset()
	ELBClassic.Reset()
	ELBV2.Reset()
	DynamoDBTables.Reset()
	ElastiCacheClusters.Reset()
	OpenSearchDomains.Reset()
	SecretsManagerSecrets.Reset()
	SFNStateMachines.Reset()
	SSMParameters.Reset()
	SQSQueues.Reset()
	SNSTopics.Reset()
	CloudFrontDistributions.Reset()
	EBSVolumes.Reset()
	VPCs.Reset()
	Subnets.Reset()
	SecurityGroups.Reset()
	NATGateways.Reset()
	InternetGateways.Reset()
	Route53HostedZones.Reset()
	ACMCertificates.Reset()
	IAMUsers.Reset()
	IAMRoles.Reset()
}
