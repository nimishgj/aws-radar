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
		[]string{"account", "account_name", "region"},
	)

	APIGatewayV2APIs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_apigatewayv2_apis_total",
			Help: "Total number of API Gateway v2 APIs",
		},
		[]string{"account", "account_name", "region", "protocol"},
	)

	// Auto Scaling Metrics
	AutoScalingGroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_groups_total",
			Help: "Total number of Auto Scaling Groups",
		},
		[]string{"account", "account_name", "region"},
	)

	// Athena Metrics
	AthenaWorkgroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_athena_workgroups_total",
			Help: "Total number of Athena workgroups",
		},
		[]string{"account", "account_name", "region"},
	)

	// ECR Metrics
	ECRRepositories = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecr_repositories_total",
			Help: "Total number of ECR repositories",
		},
		[]string{"account", "account_name", "region"},
	)

	// EC2 Metrics
	EC2Instances = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ec2_instances_total",
			Help: "Total number of EC2 instances",
		},
		[]string{"account", "account_name", "region", "instance_type", "state", "availability_zone"},
	)

	// S3 Metrics
	S3Buckets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_s3_buckets_total",
			Help: "Total number of S3 buckets",
		},
		[]string{"account", "account_name", "region"},
	)

	// RDS Metrics
	RDSInstances = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_instances_total",
			Help: "Total number of RDS instances",
		},
		[]string{"account", "account_name", "region", "db_instance_class", "engine", "multi_az", "status"},
	)

	// Lambda Metrics
	LambdaFunctions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_lambda_functions_total",
			Help: "Total number of Lambda functions",
		},
		[]string{"account", "account_name", "region", "runtime", "memory_size"},
	)

	// ECS Metrics
	ECSServices = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_services_total",
			Help: "Total number of ECS services",
		},
		[]string{"account", "account_name", "region", "cluster_name", "launch_type"},
	)

	ECSTasks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_tasks_total",
			Help: "Total number of ECS tasks",
		},
		[]string{"account", "account_name", "region", "cluster_name", "launch_type"},
	)

	ECSTasksByStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_tasks_by_status_total",
			Help: "Total number of ECS tasks by last and desired status",
		},
		[]string{"account", "account_name", "region", "cluster_name", "launch_type", "last_status", "desired_status"},
	)

	ECSClusterDepth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_cluster_depth",
			Help: "ECS cluster-level depth metrics",
		},
		[]string{"account", "account_name", "region", "cluster_name", "metric"},
	)

	// EKS Metrics
	EKSClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_eks_clusters_total",
			Help: "Total number of EKS clusters",
		},
		[]string{"account", "account_name", "region", "version", "status"},
	)

	// ELB Metrics
	ELBClassic = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elb_classic_total",
			Help: "Total number of Classic Load Balancers",
		},
		[]string{"account", "account_name", "region", "scheme"},
	)

	ELBV2 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elbv2_total",
			Help: "Total number of ALB/NLB load balancers",
		},
		[]string{"account", "account_name", "region", "type", "scheme"},
	)

	ELBV2Listeners = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elbv2_listeners_total",
			Help: "Total number of ELBv2 listeners by load balancer type and protocol",
		},
		[]string{"account", "account_name", "region", "type", "scheme", "protocol"},
	)

	ELBV2TargetGroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elbv2_target_groups_total",
			Help: "Total number of ELBv2 target groups by load balancer type and target type",
		},
		[]string{"account", "account_name", "region", "type", "target_type"},
	)

	ELBV2RulesPerALB = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elbv2_rules_per_alb",
			Help: "Number of listener rules per ALB",
		},
		[]string{"account", "account_name", "region", "load_balancer_name"},
	)

	ELBV2AvailabilityZonesPerLB = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elbv2_availability_zones_per_lb",
			Help: "Number of availability zones configured per ELBv2 load balancer",
		},
		[]string{"account", "account_name", "region", "load_balancer_name", "type", "scheme"},
	)

	ELBV2SubnetsPerLB = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elbv2_subnets_per_lb",
			Help: "Number of subnets configured per ELBv2 load balancer",
		},
		[]string{"account", "account_name", "region", "load_balancer_name", "type", "scheme"},
	)

	// DynamoDB Metrics
	DynamoDBTables = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_dynamodb_tables_total",
			Help: "Total number of DynamoDB tables",
		},
		[]string{"account", "account_name", "region", "billing_mode"},
	)

	// ElastiCache Metrics
	ElastiCacheClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elasticache_clusters_total",
			Help: "Total number of ElastiCache clusters",
		},
		[]string{"account", "account_name", "region", "engine", "cache_node_type"},
	)

	// SQS Metrics
	SQSQueues = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sqs_queues_total",
			Help: "Total number of SQS queues",
		},
		[]string{"account", "account_name", "region", "queue_type"},
	)

	// SNS Metrics
	SNSTopics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sns_topics_total",
			Help: "Total number of SNS topics",
		},
		[]string{"account", "account_name", "region"},
	)

	// CloudFront Metrics
	CloudFrontDistributions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cloudfront_distributions_total",
			Help: "Total number of CloudFront distributions",
		},
		[]string{"account", "account_name", "price_class", "enabled"},
	)

	// EBS Metrics
	EBSVolumes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ebs_volumes_total",
			Help: "Total number of EBS volumes",
		},
		[]string{"account", "account_name", "region", "volume_type", "state"},
	)

	// VPC Metrics
	VPCs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_vpc_total",
			Help: "Total number of VPCs",
		},
		[]string{"account", "account_name", "region", "state"},
	)

	Subnets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_subnet_total",
			Help: "Total number of subnets",
		},
		[]string{"account", "account_name", "region", "availability_zone"},
	)

	SecurityGroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_security_groups_total",
			Help: "Total number of security groups",
		},
		[]string{"account", "account_name", "region", "vpc_id"},
	)

	NATGateways = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_nat_gateways_total",
			Help: "Total number of NAT gateways",
		},
		[]string{"account", "account_name", "region", "state"},
	)

	InternetGateways = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_internet_gateways_total",
			Help: "Total number of Internet gateways",
		},
		[]string{"account", "account_name", "region"},
	)

	// EFS Metrics
	EFSFileSystems = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_efs_filesystems_total",
			Help: "Total number of EFS file systems",
		},
		[]string{"account", "account_name", "region"},
	)

	// EventBridge Metrics
	EventBridgeRules = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_eventbridge_rules_total",
			Help: "Total number of EventBridge rules",
		},
		[]string{"account", "account_name", "region", "event_bus"},
	)

	// Glue Metrics
	GlueJobs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_glue_jobs_total",
			Help: "Total number of Glue jobs",
		},
		[]string{"account", "account_name", "region"},
	)

	// OpenSearch Metrics
	OpenSearchDomains = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_opensearch_domains_total",
			Help: "Total number of OpenSearch domains",
		},
		[]string{"account", "account_name", "region"},
	)

	// MQ Metrics
	MQBrokers = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_mq_brokers_total",
			Help: "Total number of Amazon MQ brokers",
		},
		[]string{"account", "account_name", "region"},
	)

	// SES Metrics
	SESIdentities = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_identities_total",
			Help: "Total number of SES identities",
		},
		[]string{"account", "account_name", "region"},
	)

	// CloudFormation Metrics
	CloudFormationStacks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cloudformation_stacks_total",
			Help: "Total number of CloudFormation stacks",
		},
		[]string{"account", "account_name", "region"},
	)

	CloudFormationStacksByStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cloudformation_stacks_by_status_total",
			Help: "Total number of CloudFormation stacks by stack status",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	// DocumentDB Metrics
	DocumentDBClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_documentdb_clusters_total",
			Help: "Total number of DocumentDB clusters",
		},
		[]string{"account", "account_name", "region"},
	)

	// Neptune Metrics
	NeptuneClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_neptune_clusters_total",
			Help: "Total number of Neptune clusters",
		},
		[]string{"account", "account_name", "region"},
	)

	// MemoryDB Metrics
	MemoryDBClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_memorydb_clusters_total",
			Help: "Total number of MemoryDB clusters",
		},
		[]string{"account", "account_name", "region"},
	)

	// Timestream Metrics
	TimestreamDatabases = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_timestream_databases_total",
			Help: "Total number of Timestream databases",
		},
		[]string{"account", "account_name", "region"},
	)

	TimestreamTables = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_timestream_tables_total",
			Help: "Total number of Timestream tables",
		},
		[]string{"account", "account_name", "region"},
	)

	// FSx Metrics
	FSxFileSystems = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_fsx_filesystems_total",
			Help: "Total number of FSx file systems",
		},
		[]string{"account", "account_name", "region"},
	)

	// Backup Metrics
	BackupVaults = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_backup_vaults_total",
			Help: "Total number of AWS Backup vaults",
		},
		[]string{"account", "account_name", "region"},
	)

	// Kinesis Metrics
	KinesisStreams = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_kinesis_streams_total",
			Help: "Total number of Kinesis Data Streams",
		},
		[]string{"account", "account_name", "region"},
	)

	FirehoseDeliveryStreams = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_firehose_delivery_streams_total",
			Help: "Total number of Kinesis Data Firehose delivery streams",
		},
		[]string{"account", "account_name", "region"},
	)

	KinesisAnalyticsApplications = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_kinesisanalytics_applications_total",
			Help: "Total number of Kinesis Data Analytics applications",
		},
		[]string{"account", "account_name", "region"},
	)

	// EMR Metrics
	EMRClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_emr_clusters_total",
			Help: "Total number of EMR clusters",
		},
		[]string{"account", "account_name", "region"},
	)

	// Elastic Beanstalk Metrics
	ElasticBeanstalkApplications = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elasticbeanstalk_applications_total",
			Help: "Total number of Elastic Beanstalk applications",
		},
		[]string{"account", "account_name", "region"},
	)

	// KMS Metrics
	KMSKeys = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_kms_keys_total",
			Help: "Total number of KMS keys",
		},
		[]string{"account", "account_name", "region"},
	)

	// CloudTrail Metrics
	CloudTrailTrails = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cloudtrail_trails_total",
			Help: "Total number of CloudTrail trails",
		},
		[]string{"account", "account_name", "region"},
	)

	// Batch Metrics
	BatchJobQueues = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_batch_job_queues_total",
			Help: "Total number of AWS Batch job queues",
		},
		[]string{"account", "account_name", "region"},
	)

	// CodeBuild Metrics
	CodeBuildProjects = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_codebuild_projects_total",
			Help: "Total number of CodeBuild projects",
		},
		[]string{"account", "account_name", "region"},
	)

	// CodePipeline Metrics
	CodePipelinePipelines = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_codepipeline_pipelines_total",
			Help: "Total number of CodePipeline pipelines",
		},
		[]string{"account", "account_name", "region"},
	)

	// CodeDeploy Metrics
	CodeDeployApplications = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_codedeploy_applications_total",
			Help: "Total number of CodeDeploy applications",
		},
		[]string{"account", "account_name", "region"},
	)

	// App Runner Metrics
	AppRunnerServices = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_apprunner_services_total",
			Help: "Total number of App Runner services",
		},
		[]string{"account", "account_name", "region"},
	)

	// AWS Transfer Family Metrics
	TransferServers = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_transfer_servers_total",
			Help: "Total number of AWS Transfer Family servers",
		},
		[]string{"account", "account_name", "region"},
	)

	// MSK Metrics
	MSKClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_msk_clusters_total",
			Help: "Total number of MSK clusters",
		},
		[]string{"account", "account_name", "region"},
	)

	// Redshift Metrics
	RedshiftClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_redshift_clusters_total",
			Help: "Total number of Redshift clusters",
		},
		[]string{"account", "account_name", "region"},
	)

	// GuardDuty Metrics
	GuardDutyDetectors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_guardduty_detectors_total",
			Help: "Total number of GuardDuty detectors",
		},
		[]string{"account", "account_name", "region"},
	)

	// Security Hub Metrics
	SecurityHubStandards = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_securityhub_enabled_standards_total",
			Help: "Total number of enabled Security Hub standards",
		},
		[]string{"account", "account_name", "region"},
	)

	// Inspector2 Metrics
	InspectorCoveredResources = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_inspector2_covered_resources_total",
			Help: "Total number of Inspector2 covered resources",
		},
		[]string{"account", "account_name", "region"},
	)

	// Macie Metrics
	MacieClassificationJobs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_macie_classification_jobs_total",
			Help: "Total number of Macie classification jobs",
		},
		[]string{"account", "account_name", "region"},
	)

	// WAFv2 Metrics
	WAFWebACLs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_wafv2_web_acls_total",
			Help: "Total number of WAFv2 Web ACLs",
		},
		[]string{"account", "account_name", "region"},
	)

	// Secrets Manager Metrics
	SecretsManagerSecrets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_secretsmanager_secrets_total",
			Help: "Total number of Secrets Manager secrets",
		},
		[]string{"account", "account_name", "region"},
	)

	// SSM Metrics
	SSMParameters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ssm_parameters_total",
			Help: "Total number of SSM parameters",
		},
		[]string{"account", "account_name", "region", "type"},
	)

	// Step Functions Metrics
	SFNStateMachines = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sfn_state_machines_total",
			Help: "Total number of Step Functions state machines",
		},
		[]string{"account", "account_name", "region", "type"},
	)

	// Route53 Metrics
	Route53HostedZones = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_route53_hosted_zones_total",
			Help: "Total number of Route53 hosted zones",
		},
		[]string{"account", "account_name"},
	)

	// ACM Metrics
	ACMCertificates = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_acm_certificates_total",
			Help: "Total number of ACM certificates",
		},
		[]string{"account", "account_name", "region", "status", "type"},
	)

	// IAM Metrics
	IAMUsers = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_iam_users_total",
			Help: "Total number of IAM users",
		},
		[]string{"account", "account_name"},
	)

	IAMRoles = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_iam_roles_total",
			Help: "Total number of IAM roles",
		},
		[]string{"account", "account_name"},
	)

	// Shield Metrics
	ShieldSubscriptions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_shield_subscriptions_total",
			Help: "1 if Shield Advanced subscription exists, 0 otherwise",
		},
		[]string{"account", "account_name"},
	)

	AutoScalingGroupsWithLaunchTemplate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_groups_with_launch_template_total",
			Help: "Total number of Auto Scaling Groups configured with launch templates",
		},
		[]string{"account", "account_name", "region"},
	)

	AutoScalingPolicies = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_policies_total",
			Help: "Total number of Auto Scaling policies",
		},
		[]string{"account", "account_name", "region"},
	)

	AutoScalingPoliciesByType = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_policies_by_type_total",
			Help: "Total number of Auto Scaling policies by policy type",
		},
		[]string{"account", "account_name", "region", "policy_type"},
	)

	AutoScalingGroupsByMixedInstances = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_groups_by_mixed_instances_total",
			Help: "Total number of Auto Scaling groups by mixed instances policy presence",
		},
		[]string{"account", "account_name", "region", "has_mixed_instances_policy"},
	)

	AutoScalingLaunchTemplateUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_launch_template_usage_total",
			Help: "Total number of Auto Scaling groups by launch template and version",
		},
		[]string{"account", "account_name", "region", "launch_template_id", "launch_template_version"},
	)

	AutoScalingLifecycleHooks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_lifecycle_hooks_total",
			Help: "Total number of Auto Scaling lifecycle hooks",
		},
		[]string{"account", "account_name", "region"},
	)

	AutoScalingWarmPools = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_warm_pools_total",
			Help: "Total number of Auto Scaling groups with warm pools",
		},
		[]string{"account", "account_name", "region"},
	)

	AutoScalingWarmPoolInstances = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_warm_pool_instances_total",
			Help: "Total number of Auto Scaling warm pool instances",
		},
		[]string{"account", "account_name", "region"},
	)

	AutoScalingInstanceRefreshes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_autoscaling_instance_refreshes_total",
			Help: "Total number of Auto Scaling instance refreshes by status",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	ELBV2Detailed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elbv2_detailed_total",
			Help: "Total number of ALB/NLB by type, scheme, IP type and state",
		},
		[]string{"account", "account_name", "region", "type", "scheme", "ip_address_type", "state"},
	)

	ECSServicesByStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_services_by_status_total",
			Help: "Total number of ECS services by status",
		},
		[]string{"account", "account_name", "region", "cluster_name", "launch_type", "status"},
	)

	ECSCapacityProviders = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_capacity_providers_total",
			Help: "Total number of ECS capacity providers by status",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	ECSCapacityProvidersDetailed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_capacity_providers_detailed_total",
			Help: "Total number of ECS capacity providers by type and status",
		},
		[]string{"account", "account_name", "region", "capacity_provider_type", "status"},
	)

	ECSDefaultCapacityProviderStrategy = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_default_capacity_provider_strategy_total",
			Help: "Number of ECS clusters using each capacity provider in default strategy",
		},
		[]string{"account", "account_name", "region", "cluster_name", "capacity_provider"},
	)

	ECSTaskDefinitions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_task_definitions_total",
			Help: "Total number of ECS task definitions by status",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	ECSTaskDefinitionsDetailed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_task_definitions_detailed_total",
			Help: "Total number of ECS task definitions by family, revision and runtime platform",
		},
		[]string{"account", "account_name", "region", "status", "family", "revision", "os_family", "cpu_architecture"},
	)

	ECRPublicRepositories = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecr_public_repositories_total",
			Help: "Total number of ECR Public repositories",
		},
		[]string{"account", "account_name"},
	)

	Route53ResolverEndpoints = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_route53resolver_endpoints_total",
			Help: "Total number of Route53 Resolver endpoints",
		},
		[]string{"account", "account_name", "region", "direction", "status"},
	)

	Route53ResolverRules = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_route53resolver_rules_total",
			Help: "Total number of Route53 Resolver rules",
		},
		[]string{"account", "account_name", "region", "rule_type"},
	)

	VPCEndpoints = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_vpc_endpoints_total",
			Help: "Total number of VPC endpoints",
		},
		[]string{"account", "account_name", "region", "endpoint_type", "state"},
	)

	TransitGateways = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_transit_gateways_total",
			Help: "Total number of Transit Gateways",
		},
		[]string{"account", "account_name", "region", "state"},
	)

	VPNGateways = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_vpn_gateways_total",
			Help: "Total number of VPN Gateways",
		},
		[]string{"account", "account_name", "region", "state"},
	)

	OrganizationsAccounts = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_organizations_accounts_total",
			Help: "Total number of AWS Organizations accounts",
		},
		[]string{"account", "account_name", "state"},
	)

	OrganizationsOrganizationalUnits = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_organizations_organizational_units_total",
			Help: "Total number of AWS Organizations OUs",
		},
		[]string{"account", "account_name"},
	)

	ControlTowerLandingZones = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_controltower_landing_zones_total",
			Help: "Total number of AWS Control Tower landing zones",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	ConfigRecorders = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_config_recorders_total",
			Help: "Total number of AWS Config configuration recorders",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	CloudTrailLakeEventDataStores = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cloudtrail_lake_event_data_stores_total",
			Help: "Total number of CloudTrail Lake event data stores",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	KinesisVideoStreams = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_kinesisvideo_streams_total",
			Help: "Total number of Kinesis Video Streams",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	OpenSearchServerlessCollections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_opensearch_serverless_collections_total",
			Help: "Total number of OpenSearch Serverless collections",
		},
		[]string{"account", "account_name", "region", "status", "type"},
	)

	SESConfigSets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_config_sets_total",
			Help: "Total number of SES configuration sets",
		},
		[]string{"account", "account_name", "region"},
	)

	SESContactLists = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_contact_lists_total",
			Help: "Total number of SES contact lists",
		},
		[]string{"account", "account_name", "region"},
	)

	SESIdentitiesByVerificationStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_identities_by_verification_status_total",
			Help: "Total number of SES identities by verification status",
		},
		[]string{"account", "account_name", "region", "verification_status"},
	)

	SESIdentityAuthStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_identity_auth_status_total",
			Help: "Total number of SES identities by DKIM, SPF and MAIL FROM status",
		},
		[]string{"account", "account_name", "region", "dkim_status", "spf_status", "mail_from_status"},
	)

	SESConfigSetEventDestinations = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_config_set_event_destinations_total",
			Help: "Total number of SES configuration set event destinations",
		},
		[]string{"account", "account_name", "region", "event_destination_type"},
	)

	SESSuppressedDestinations = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_suppressed_destinations_total",
			Help: "Total number of SES suppressed destinations by reason",
		},
		[]string{"account", "account_name", "region", "reason"},
	)

	SESDedicatedIPPools = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_dedicated_ip_pools_total",
			Help: "Total number of SES dedicated IP pools",
		},
		[]string{"account", "account_name", "region"},
	)

	SESConfigSetsBySendingPool = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_config_sets_by_sending_pool_total",
			Help: "Total number of SES configuration sets by dedicated sending pool",
		},
		[]string{"account", "account_name", "region", "sending_pool"},
	)

	SESAccountSettings = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_account_settings",
			Help: "SES account setting values",
		},
		[]string{"account", "account_name", "region", "setting"},
	)

	SESSendingQuota = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ses_sending_quota",
			Help: "SES sending quota values",
		},
		[]string{"account", "account_name", "region", "quota_type"},
	)

	S3AccessPoints = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_s3_access_points_total",
			Help: "Total number of S3 Access Points",
		},
		[]string{"account", "account_name", "region"},
	)

	S3StorageLensConfigurations = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_s3_storage_lens_configurations_total",
			Help: "Total number of S3 Storage Lens configurations",
		},
		[]string{"account", "account_name", "region"},
	)

	RDSProxies = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_proxies_total",
			Help: "Total number of RDS proxies",
		},
		[]string{"account", "account_name", "region", "engine_family", "status"},
	)

	RDSAuroraServerlessClusters = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_aurora_serverless_clusters_total",
			Help: "Total number of Aurora serverless clusters",
		},
		[]string{"account", "account_name", "region", "engine", "engine_mode"},
	)

	RDSInstancesByEngineVersion = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_instances_by_engine_version_total",
			Help: "Total number of RDS instances by engine and engine version",
		},
		[]string{"account", "account_name", "region", "engine", "engine_version"},
	)

	RDSInstancesByClass = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_instances_by_class_total",
			Help: "Total number of RDS instances by DB instance class",
		},
		[]string{"account", "account_name", "region", "db_instance_class"},
	)

	RDSInstancesByMultiAZ = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_instances_by_multi_az_total",
			Help: "Total number of RDS instances by Multi-AZ setting",
		},
		[]string{"account", "account_name", "region", "multi_az"},
	)

	RDSReadReplicas = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_read_replicas_total",
			Help: "Total number of RDS read replica instances by engine",
		},
		[]string{"account", "account_name", "region", "engine"},
	)

	RDSProxyTargets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_proxy_targets_total",
			Help: "Total number of RDS proxy targets by engine family and target type",
		},
		[]string{"account", "account_name", "region", "engine_family", "target_type"},
	)

	RDSAuroraServerlessV2Capacity = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_aurora_serverless_v2_capacity",
			Help: "Aurora Serverless v2 min and max capacity values",
		},
		[]string{"account", "account_name", "region", "cluster_identifier", "metric"},
	)

	RDSAuroraServerlessByStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_rds_aurora_serverless_by_status_total",
			Help: "Total number of Aurora serverless clusters by status",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	ElastiCacheServerlessCaches = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elasticache_serverless_caches_total",
			Help: "Total number of ElastiCache serverless caches",
		},
		[]string{"account", "account_name", "region", "engine", "status"},
	)

	ElastiCacheReplicationGroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elasticache_replication_groups_total",
			Help: "Total number of ElastiCache replication groups",
		},
		[]string{"account", "account_name", "region", "engine", "status", "cluster_enabled"},
	)

	ElastiCacheGlobalReplicationGroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elasticache_global_replication_groups_total",
			Help: "Total number of ElastiCache global replication groups by status",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	BedrockCustomModels = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_bedrock_custom_models_total",
			Help: "Total number of Bedrock custom models",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	SageMakerEndpoints = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sagemaker_endpoints_total",
			Help: "Total number of SageMaker endpoints",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	QuickSightDashboards = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_quicksight_dashboards_total",
			Help: "Total number of QuickSight dashboards",
		},
		[]string{"account", "account_name", "region"},
	)

	WorkSpacesInstances = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_workspaces_instances_total",
			Help: "Total number of WorkSpaces instances",
		},
		[]string{"account", "account_name", "region", "state"},
	)

	AppStreamFleets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_appstream_fleets_total",
			Help: "Total number of AppStream fleets",
		},
		[]string{"account", "account_name", "region", "state"},
	)

	ConnectInstances = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_connect_instances_total",
			Help: "Total number of Amazon Connect instances",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	AmplifyApps = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_amplify_apps_total",
			Help: "Total number of Amplify apps",
		},
		[]string{"account", "account_name", "region"},
	)

	GlobalAccelerators = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_globalaccelerator_accelerators_total",
			Help: "Total number of AWS Global Accelerator accelerators",
		},
		[]string{"account", "account_name", "ip_address_type", "enabled"},
	)

	DataSyncTasks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_datasync_tasks_total",
			Help: "Total number of AWS DataSync tasks",
		},
		[]string{"account", "account_name", "region"},
	)

	DMSReplicationInstances = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_dms_replication_instances_total",
			Help: "Total number of AWS DMS replication instances",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	DirectConnectConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_directconnect_connections_total",
			Help: "Total number of Direct Connect connections",
		},
		[]string{"account", "account_name", "state"},
	)

	// Pricing Metrics
	CostByService = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cost_service_unblended_usd",
			Help: "Daily unblended cost by AWS service in USD",
		},
		[]string{"account", "account_name", "service", "period_start"},
	)

	CostTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cost_total_unblended_usd",
			Help: "Daily total unblended AWS cost in USD",
		},
		[]string{"account", "account_name", "period_start"},
	)

	// CUR (Cost and Usage Report) Metrics
	CURTotalCost = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cur_total_cost_usd",
			Help: "Total cost from CUR report in USD",
		},
		[]string{"account", "account_name", "period"},
	)

	CURCostByService = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cur_cost_by_service_usd",
			Help: "Cost by AWS service from CUR report in USD",
		},
		[]string{"account", "account_name", "service", "period"},
	)

	CURCostByResource = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cur_cost_by_resource_usd",
			Help: "Cost by resource ID from CUR report in USD (top N)",
		},
		[]string{"account", "account_name", "service", "resource_id", "period"},
	)

	CURCostByUsageType = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cur_cost_by_usage_type_usd",
			Help: "Cost by usage type from CUR report in USD",
		},
		[]string{"account", "account_name", "service", "usage_type", "period"},
	)

	CURCostByTag = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cur_cost_by_tag_usd",
			Help: "Cost by user tag from CUR report in USD",
		},
		[]string{"account", "account_name", "tag_key", "tag_value", "period"},
	)

	CURLastProcessed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cur_last_processed_timestamp",
			Help: "Unix timestamp of last successful CUR processing",
		},
		[]string{"account", "account_name"},
	)

	// Cognito Metrics
	CognitoUserPools = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cognito_user_pools_total",
			Help: "Total number of Cognito user pools",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	CognitoUserPoolUsers = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cognito_user_pool_users_total",
			Help: "Estimated number of users per Cognito user pool",
		},
		[]string{"account", "account_name", "region", "user_pool_name"},
	)

	CognitoIdentityPools = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cognito_identity_pools_total",
			Help: "Total number of Cognito identity pools",
		},
		[]string{"account", "account_name", "region"},
	)

	// Network Firewall Metrics
	NetworkFirewallFirewalls = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_networkfirewall_firewalls_total",
			Help: "Total number of AWS Network Firewall firewalls",
		},
		[]string{"account", "account_name", "region"},
	)

	NetworkFirewallPolicies = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_networkfirewall_policies_total",
			Help: "Total number of AWS Network Firewall policies",
		},
		[]string{"account", "account_name", "region"},
	)

	NetworkFirewallRuleGroups = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_networkfirewall_rule_groups_total",
			Help: "Total number of AWS Network Firewall rule groups",
		},
		[]string{"account", "account_name", "region", "type"},
	)

	// Firewall Manager Metrics
	FMSPolicies = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_fms_policies_total",
			Help: "Total number of Firewall Manager policies",
		},
		[]string{"account", "account_name", "region", "resource_type"},
	)

	// Private CA Metrics
	ACMPCACertificateAuthorities = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_acmpca_certificate_authorities_total",
			Help: "Total number of ACM Private Certificate Authorities",
		},
		[]string{"account", "account_name", "region", "status", "type"},
	)

	// Service Catalog Metrics
	ServiceCatalogPortfolios = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_servicecatalog_portfolios_total",
			Help: "Total number of Service Catalog portfolios",
		},
		[]string{"account", "account_name", "region"},
	)

	ServiceCatalogProducts = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_servicecatalog_products_total",
			Help: "Total number of Service Catalog products",
		},
		[]string{"account", "account_name", "region"},
	)

	ServiceCatalogProvisionedProducts = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_servicecatalog_provisioned_products_total",
			Help: "Total number of Service Catalog provisioned products",
		},
		[]string{"account", "account_name", "region", "status"},
	)

	// License Manager Metrics
	LicenseManagerConfigurations = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_licensemanager_license_configurations_total",
			Help: "Total number of License Manager configurations",
		},
		[]string{"account", "account_name", "region"},
	)

	// SNS Depth Metrics
	SNSSubscriptions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sns_subscriptions_total",
			Help: "Total number of confirmed SNS subscriptions",
		},
		[]string{"account", "account_name", "region"},
	)

	SNSTopicsByType = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sns_topics_by_type_total",
			Help: "Total number of SNS topics by type (FIFO vs standard)",
		},
		[]string{"account", "account_name", "region", "fifo"},
	)

	// SQS Depth Metrics
	SQSMessages = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sqs_messages_total",
			Help: "Total approximate number of messages across SQS queues",
		},
		[]string{"account", "account_name", "region", "queue_type"},
	)

	SQSQueuesWithDLQ = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_sqs_queues_with_dlq_total",
			Help: "Total number of SQS queues with dead-letter queue configured",
		},
		[]string{"account", "account_name", "region"},
	)

	// SSM Depth Metrics
	SSMDocuments = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ssm_documents_total",
			Help: "Total number of SSM documents",
		},
		[]string{"account", "account_name", "region", "owner"},
	)

	SSMMaintenanceWindows = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ssm_maintenance_windows_total",
			Help: "Total number of SSM maintenance windows",
		},
		[]string{"account", "account_name", "region", "enabled"},
	)

	SSMAssociations = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ssm_associations_total",
			Help: "Total number of SSM associations",
		},
		[]string{"account", "account_name", "region"},
	)

	SSMPatchBaselines = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ssm_patch_baselines_total",
			Help: "Total number of SSM patch baselines",
		},
		[]string{"account", "account_name", "region"},
	)

	// Collection Metrics
	CollectionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "aws_radar_collection_duration_seconds",
			Help:    "Duration of AWS resource collection",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"account", "account_name", "collector"},
	)

	CollectionUp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_radar_up",
			Help: "1 if the last collection succeeded, 0 otherwise",
		},
		[]string{"account", "account_name", "collector", "region"},
	)

	CollectionErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aws_radar_collection_errors_total",
			Help: "Total number of collection errors",
		},
		[]string{"account", "account_name", "collector", "region"},
	)
)

// ResetAll resets all gauge metrics before a new collection cycle
func ResetAll() {
	APIGatewayRestAPIs.Reset()
	APIGatewayV2APIs.Reset()
	AutoScalingGroups.Reset()
	AutoScalingGroupsWithLaunchTemplate.Reset()
	AutoScalingPolicies.Reset()
	AutoScalingPoliciesByType.Reset()
	AutoScalingGroupsByMixedInstances.Reset()
	AutoScalingLaunchTemplateUsage.Reset()
	AutoScalingLifecycleHooks.Reset()
	AutoScalingWarmPools.Reset()
	AutoScalingWarmPoolInstances.Reset()
	AutoScalingInstanceRefreshes.Reset()
	AthenaWorkgroups.Reset()
	ECRRepositories.Reset()
	ECRPublicRepositories.Reset()
	EC2Instances.Reset()
	EFSFileSystems.Reset()
	EventBridgeRules.Reset()
	GlueJobs.Reset()
	S3Buckets.Reset()
	S3AccessPoints.Reset()
	S3StorageLensConfigurations.Reset()
	RDSInstances.Reset()
	RDSProxies.Reset()
	RDSAuroraServerlessClusters.Reset()
	RDSInstancesByEngineVersion.Reset()
	RDSInstancesByClass.Reset()
	RDSInstancesByMultiAZ.Reset()
	RDSReadReplicas.Reset()
	RDSProxyTargets.Reset()
	RDSAuroraServerlessV2Capacity.Reset()
	RDSAuroraServerlessByStatus.Reset()
	LambdaFunctions.Reset()
	ECSServices.Reset()
	ECSTasks.Reset()
	ECSTasksByStatus.Reset()
	ECSClusterDepth.Reset()
	ECSServicesByStatus.Reset()
	ECSCapacityProviders.Reset()
	ECSCapacityProvidersDetailed.Reset()
	ECSDefaultCapacityProviderStrategy.Reset()
	ECSTaskDefinitions.Reset()
	ECSTaskDefinitionsDetailed.Reset()
	EKSClusters.Reset()
	ELBClassic.Reset()
	ELBV2.Reset()
	ELBV2Detailed.Reset()
	ELBV2Listeners.Reset()
	ELBV2TargetGroups.Reset()
	ELBV2RulesPerALB.Reset()
	ELBV2AvailabilityZonesPerLB.Reset()
	ELBV2SubnetsPerLB.Reset()
	DynamoDBTables.Reset()
	ElastiCacheClusters.Reset()
	ElastiCacheServerlessCaches.Reset()
	ElastiCacheReplicationGroups.Reset()
	ElastiCacheGlobalReplicationGroups.Reset()
	OpenSearchDomains.Reset()
	OpenSearchServerlessCollections.Reset()
	MQBrokers.Reset()
	SESIdentities.Reset()
	SESConfigSets.Reset()
	SESContactLists.Reset()
	SESIdentitiesByVerificationStatus.Reset()
	SESIdentityAuthStatus.Reset()
	SESConfigSetEventDestinations.Reset()
	SESSuppressedDestinations.Reset()
	SESDedicatedIPPools.Reset()
	SESConfigSetsBySendingPool.Reset()
	SESAccountSettings.Reset()
	SESSendingQuota.Reset()
	CloudFormationStacks.Reset()
	CloudFormationStacksByStatus.Reset()
	DocumentDBClusters.Reset()
	NeptuneClusters.Reset()
	MemoryDBClusters.Reset()
	TimestreamDatabases.Reset()
	TimestreamTables.Reset()
	FSxFileSystems.Reset()
	BackupVaults.Reset()
	KinesisStreams.Reset()
	FirehoseDeliveryStreams.Reset()
	KinesisAnalyticsApplications.Reset()
	KinesisVideoStreams.Reset()
	EMRClusters.Reset()
	ElasticBeanstalkApplications.Reset()
	KMSKeys.Reset()
	CloudTrailTrails.Reset()
	CloudTrailLakeEventDataStores.Reset()
	BatchJobQueues.Reset()
	CodeBuildProjects.Reset()
	CodePipelinePipelines.Reset()
	CodeDeployApplications.Reset()
	AppRunnerServices.Reset()
	TransferServers.Reset()
	MSKClusters.Reset()
	RedshiftClusters.Reset()
	GuardDutyDetectors.Reset()
	SecurityHubStandards.Reset()
	InspectorCoveredResources.Reset()
	MacieClassificationJobs.Reset()
	WAFWebACLs.Reset()
	SecretsManagerSecrets.Reset()
	SFNStateMachines.Reset()
	SSMParameters.Reset()
	CollectionUp.Reset()
	SQSQueues.Reset()
	SNSTopics.Reset()
	CloudFrontDistributions.Reset()
	EBSVolumes.Reset()
	VPCs.Reset()
	Subnets.Reset()
	SecurityGroups.Reset()
	NATGateways.Reset()
	InternetGateways.Reset()
	VPCEndpoints.Reset()
	TransitGateways.Reset()
	VPNGateways.Reset()
	Route53HostedZones.Reset()
	Route53ResolverEndpoints.Reset()
	Route53ResolverRules.Reset()
	ACMCertificates.Reset()
	IAMUsers.Reset()
	IAMRoles.Reset()
	ShieldSubscriptions.Reset()
	OrganizationsAccounts.Reset()
	OrganizationsOrganizationalUnits.Reset()
	ControlTowerLandingZones.Reset()
	ConfigRecorders.Reset()
	DirectConnectConnections.Reset()
	BedrockCustomModels.Reset()
	SageMakerEndpoints.Reset()
	QuickSightDashboards.Reset()
	WorkSpacesInstances.Reset()
	AppStreamFleets.Reset()
	ConnectInstances.Reset()
	AmplifyApps.Reset()
	GlobalAccelerators.Reset()
	DataSyncTasks.Reset()
	DMSReplicationInstances.Reset()
	CostByService.Reset()
	CostTotal.Reset()
	CURTotalCost.Reset()
	CURCostByService.Reset()
	CURCostByResource.Reset()
	CURCostByUsageType.Reset()
	CURCostByTag.Reset()
	CURLastProcessed.Reset()
	CognitoUserPools.Reset()
	CognitoUserPoolUsers.Reset()
	CognitoIdentityPools.Reset()
	NetworkFirewallFirewalls.Reset()
	NetworkFirewallPolicies.Reset()
	NetworkFirewallRuleGroups.Reset()
	FMSPolicies.Reset()
	ACMPCACertificateAuthorities.Reset()
	ServiceCatalogPortfolios.Reset()
	ServiceCatalogProducts.Reset()
	ServiceCatalogProvisionedProducts.Reset()
	LicenseManagerConfigurations.Reset()
	SNSSubscriptions.Reset()
	SNSTopicsByType.Reset()
	SQSMessages.Reset()
	SQSQueuesWithDLQ.Reset()
	SSMDocuments.Reset()
	SSMMaintenanceWindows.Reset()
	SSMAssociations.Reset()
	SSMPatchBaselines.Reset()
}

// InitRegionalDefaults ensures every regional metric has at least one 0-valued
// time series so that Grafana panels show "0" instead of "no data" when there
// are no resources or a collector errors out.
func InitRegionalDefaults(account, accountName, region string) {
	n := "none" // sentinel label for dimensional metrics

	// Simple counters (account, account_name, region only)
	APIGatewayRestAPIs.WithLabelValues(account, accountName, region).Add(0)
	AthenaWorkgroups.WithLabelValues(account, accountName, region).Add(0)
	AutoScalingGroups.WithLabelValues(account, accountName, region).Add(0)
	AutoScalingGroupsWithLaunchTemplate.WithLabelValues(account, accountName, region).Add(0)
	AutoScalingPolicies.WithLabelValues(account, accountName, region).Add(0)
	AutoScalingLifecycleHooks.WithLabelValues(account, accountName, region).Add(0)
	AutoScalingWarmPools.WithLabelValues(account, accountName, region).Add(0)
	AutoScalingWarmPoolInstances.WithLabelValues(account, accountName, region).Add(0)
	BackupVaults.WithLabelValues(account, accountName, region).Add(0)
	BatchJobQueues.WithLabelValues(account, accountName, region).Add(0)
	CodeBuildProjects.WithLabelValues(account, accountName, region).Add(0)
	CodePipelinePipelines.WithLabelValues(account, accountName, region).Add(0)
	CodeDeployApplications.WithLabelValues(account, accountName, region).Add(0)
	AppRunnerServices.WithLabelValues(account, accountName, region).Add(0)
	TransferServers.WithLabelValues(account, accountName, region).Add(0)
	MSKClusters.WithLabelValues(account, accountName, region).Add(0)
	RedshiftClusters.WithLabelValues(account, accountName, region).Add(0)
	GuardDutyDetectors.WithLabelValues(account, accountName, region).Add(0)
	SecurityHubStandards.WithLabelValues(account, accountName, region).Add(0)
	InspectorCoveredResources.WithLabelValues(account, accountName, region).Add(0)
	MacieClassificationJobs.WithLabelValues(account, accountName, region).Add(0)
	WAFWebACLs.WithLabelValues(account, accountName, region).Add(0)
	SecretsManagerSecrets.WithLabelValues(account, accountName, region).Add(0)
	CloudFormationStacks.WithLabelValues(account, accountName, region).Add(0)
	DocumentDBClusters.WithLabelValues(account, accountName, region).Add(0)
	NeptuneClusters.WithLabelValues(account, accountName, region).Add(0)
	MemoryDBClusters.WithLabelValues(account, accountName, region).Add(0)
	TimestreamDatabases.WithLabelValues(account, accountName, region).Add(0)
	TimestreamTables.WithLabelValues(account, accountName, region).Add(0)
	FSxFileSystems.WithLabelValues(account, accountName, region).Add(0)
	KinesisStreams.WithLabelValues(account, accountName, region).Add(0)
	FirehoseDeliveryStreams.WithLabelValues(account, accountName, region).Add(0)
	KinesisAnalyticsApplications.WithLabelValues(account, accountName, region).Add(0)
	EMRClusters.WithLabelValues(account, accountName, region).Add(0)
	ElasticBeanstalkApplications.WithLabelValues(account, accountName, region).Add(0)
	KMSKeys.WithLabelValues(account, accountName, region).Add(0)
	CloudTrailTrails.WithLabelValues(account, accountName, region).Add(0)
	OpenSearchDomains.WithLabelValues(account, accountName, region).Add(0)
	MQBrokers.WithLabelValues(account, accountName, region).Add(0)
	SESIdentities.WithLabelValues(account, accountName, region).Add(0)
	SESConfigSets.WithLabelValues(account, accountName, region).Add(0)
	SESContactLists.WithLabelValues(account, accountName, region).Add(0)
	SESDedicatedIPPools.WithLabelValues(account, accountName, region).Add(0)
	EFSFileSystems.WithLabelValues(account, accountName, region).Add(0)
	GlueJobs.WithLabelValues(account, accountName, region).Add(0)
	ECRRepositories.WithLabelValues(account, accountName, region).Add(0)
	InternetGateways.WithLabelValues(account, accountName, region).Add(0)
	AmplifyApps.WithLabelValues(account, accountName, region).Add(0)
	DataSyncTasks.WithLabelValues(account, accountName, region).Add(0)
	QuickSightDashboards.WithLabelValues(account, accountName, region).Add(0)
	SNSTopics.WithLabelValues(account, accountName, region).Add(0)
	SNSSubscriptions.WithLabelValues(account, accountName, region).Add(0)
	SQSQueuesWithDLQ.WithLabelValues(account, accountName, region).Add(0)
	SSMAssociations.WithLabelValues(account, accountName, region).Add(0)
	SSMPatchBaselines.WithLabelValues(account, accountName, region).Add(0)
	NetworkFirewallFirewalls.WithLabelValues(account, accountName, region).Add(0)
	NetworkFirewallPolicies.WithLabelValues(account, accountName, region).Add(0)
	ServiceCatalogPortfolios.WithLabelValues(account, accountName, region).Add(0)
	ServiceCatalogProducts.WithLabelValues(account, accountName, region).Add(0)
	LicenseManagerConfigurations.WithLabelValues(account, accountName, region).Add(0)
	CognitoIdentityPools.WithLabelValues(account, accountName, region).Add(0)
	S3Buckets.WithLabelValues(account, accountName, region).Add(0)
	S3AccessPoints.WithLabelValues(account, accountName, region).Add(0)
	S3StorageLensConfigurations.WithLabelValues(account, accountName, region).Add(0)

	// Dimensional metrics — use "none" so panels show 0 when no resources exist
	APIGatewayV2APIs.WithLabelValues(account, accountName, region, n).Add(0)
	EC2Instances.WithLabelValues(account, accountName, region, n, n, n).Add(0)
	RDSInstances.WithLabelValues(account, accountName, region, n, n, n, n).Add(0)
	RDSProxies.WithLabelValues(account, accountName, region, n, n).Add(0)
	RDSAuroraServerlessClusters.WithLabelValues(account, accountName, region, n, n).Add(0)
	RDSInstancesByEngineVersion.WithLabelValues(account, accountName, region, n, n).Add(0)
	RDSInstancesByClass.WithLabelValues(account, accountName, region, n).Add(0)
	RDSInstancesByMultiAZ.WithLabelValues(account, accountName, region, n).Add(0)
	RDSReadReplicas.WithLabelValues(account, accountName, region, n).Add(0)
	RDSProxyTargets.WithLabelValues(account, accountName, region, n, n).Add(0)
	RDSAuroraServerlessByStatus.WithLabelValues(account, accountName, region, n).Add(0)
	LambdaFunctions.WithLabelValues(account, accountName, region, n, n).Add(0)
	ECSServices.WithLabelValues(account, accountName, region, n, n).Add(0)
	ECSTasks.WithLabelValues(account, accountName, region, n, n).Add(0)
	ECSTasksByStatus.WithLabelValues(account, accountName, region, n, n, n, n).Add(0)
	ECSServicesByStatus.WithLabelValues(account, accountName, region, n, n, n).Add(0)
	ECSCapacityProviders.WithLabelValues(account, accountName, region, n).Add(0)
	ECSCapacityProvidersDetailed.WithLabelValues(account, accountName, region, n, n).Add(0)
	ECSTaskDefinitions.WithLabelValues(account, accountName, region, n).Add(0)
	EKSClusters.WithLabelValues(account, accountName, region, n, n).Add(0)
	ELBClassic.WithLabelValues(account, accountName, region, n).Add(0)
	ELBV2.WithLabelValues(account, accountName, region, n, n).Add(0)
	ELBV2Detailed.WithLabelValues(account, accountName, region, n, n, n, n).Add(0)
	ELBV2Listeners.WithLabelValues(account, accountName, region, n, n, n).Add(0)
	ELBV2TargetGroups.WithLabelValues(account, accountName, region, n, n).Add(0)
	DynamoDBTables.WithLabelValues(account, accountName, region, n).Add(0)
	ElastiCacheClusters.WithLabelValues(account, accountName, region, n, n).Add(0)
	ElastiCacheServerlessCaches.WithLabelValues(account, accountName, region, n, n).Add(0)
	ElastiCacheReplicationGroups.WithLabelValues(account, accountName, region, n, n, n).Add(0)
	ElastiCacheGlobalReplicationGroups.WithLabelValues(account, accountName, region, n).Add(0)
	SQSQueues.WithLabelValues(account, accountName, region, n).Add(0)
	SQSMessages.WithLabelValues(account, accountName, region, n).Add(0)
	EBSVolumes.WithLabelValues(account, accountName, region, n, n).Add(0)
	VPCs.WithLabelValues(account, accountName, region, n).Add(0)
	Subnets.WithLabelValues(account, accountName, region, n).Add(0)
	SecurityGroups.WithLabelValues(account, accountName, region, n).Add(0)
	NATGateways.WithLabelValues(account, accountName, region, n).Add(0)
	VPCEndpoints.WithLabelValues(account, accountName, region, n, n).Add(0)
	TransitGateways.WithLabelValues(account, accountName, region, n).Add(0)
	VPNGateways.WithLabelValues(account, accountName, region, n).Add(0)
	Route53ResolverEndpoints.WithLabelValues(account, accountName, region, n, n).Add(0)
	Route53ResolverRules.WithLabelValues(account, accountName, region, n).Add(0)
	ACMCertificates.WithLabelValues(account, accountName, region, n, n).Add(0)
	ControlTowerLandingZones.WithLabelValues(account, accountName, region, n).Add(0)
	ConfigRecorders.WithLabelValues(account, accountName, region, n).Add(0)
	CloudTrailLakeEventDataStores.WithLabelValues(account, accountName, region, n).Add(0)
	KinesisVideoStreams.WithLabelValues(account, accountName, region, n).Add(0)
	OpenSearchServerlessCollections.WithLabelValues(account, accountName, region, n, n).Add(0)
	CloudFormationStacksByStatus.WithLabelValues(account, accountName, region, n).Add(0)
	EventBridgeRules.WithLabelValues(account, accountName, region, n).Add(0)
	SSMParameters.WithLabelValues(account, accountName, region, n).Add(0)
	SFNStateMachines.WithLabelValues(account, accountName, region, n).Add(0)
	AutoScalingPoliciesByType.WithLabelValues(account, accountName, region, n).Add(0)
	AutoScalingGroupsByMixedInstances.WithLabelValues(account, accountName, region, n).Add(0)
	AutoScalingInstanceRefreshes.WithLabelValues(account, accountName, region, n).Add(0)
	SESIdentitiesByVerificationStatus.WithLabelValues(account, accountName, region, n).Add(0)
	SESIdentityAuthStatus.WithLabelValues(account, accountName, region, n, n, n).Add(0)
	SESConfigSetEventDestinations.WithLabelValues(account, accountName, region, n).Add(0)
	SESSuppressedDestinations.WithLabelValues(account, accountName, region, n).Add(0)
	SESConfigSetsBySendingPool.WithLabelValues(account, accountName, region, n).Add(0)
	SESAccountSettings.WithLabelValues(account, accountName, region, n).Add(0)
	SESSendingQuota.WithLabelValues(account, accountName, region, n).Add(0)
	BedrockCustomModels.WithLabelValues(account, accountName, region, n).Add(0)
	SageMakerEndpoints.WithLabelValues(account, accountName, region, n).Add(0)
	WorkSpacesInstances.WithLabelValues(account, accountName, region, n).Add(0)
	AppStreamFleets.WithLabelValues(account, accountName, region, n).Add(0)
	ConnectInstances.WithLabelValues(account, accountName, region, n).Add(0)
	DMSReplicationInstances.WithLabelValues(account, accountName, region, n).Add(0)
	NetworkFirewallRuleGroups.WithLabelValues(account, accountName, region, n).Add(0)
	FMSPolicies.WithLabelValues(account, accountName, region, n).Add(0)
	ACMPCACertificateAuthorities.WithLabelValues(account, accountName, region, n, n).Add(0)
	ServiceCatalogProvisionedProducts.WithLabelValues(account, accountName, region, n).Add(0)
	CognitoUserPools.WithLabelValues(account, accountName, region, n).Add(0)
	CognitoUserPoolUsers.WithLabelValues(account, accountName, region, n).Add(0)
	SNSTopicsByType.WithLabelValues(account, accountName, region, n).Add(0)
	SSMDocuments.WithLabelValues(account, accountName, region, n).Add(0)
	SSMMaintenanceWindows.WithLabelValues(account, accountName, region, n).Add(0)
}

// InitGlobalDefaults ensures every global metric has at least one 0-valued
// time series.
func InitGlobalDefaults(account, accountName string) {
	n := "none"

	Route53HostedZones.WithLabelValues(account, accountName).Add(0)
	IAMUsers.WithLabelValues(account, accountName).Add(0)
	IAMRoles.WithLabelValues(account, accountName).Add(0)
	ShieldSubscriptions.WithLabelValues(account, accountName).Add(0)
	ECRPublicRepositories.WithLabelValues(account, accountName).Add(0)
	OrganizationsOrganizationalUnits.WithLabelValues(account, accountName).Add(0)
	CloudFrontDistributions.WithLabelValues(account, accountName, n, n).Add(0)
	OrganizationsAccounts.WithLabelValues(account, accountName, n).Add(0)
	DirectConnectConnections.WithLabelValues(account, accountName, n).Add(0)
	GlobalAccelerators.WithLabelValues(account, accountName, n, n).Add(0)
	CURTotalCost.WithLabelValues(account, accountName, n).Add(0)
	CURCostByService.WithLabelValues(account, accountName, n, n).Add(0)
	CURLastProcessed.WithLabelValues(account, accountName).Add(0)
}
