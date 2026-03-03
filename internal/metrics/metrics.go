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

	ECSTaskDefinitions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_ecs_task_definitions_total",
			Help: "Total number of ECS task definitions by status",
		},
		[]string{"account", "account_name", "region", "status"},
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

	ElastiCacheServerlessCaches = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_elasticache_serverless_caches_total",
			Help: "Total number of ElastiCache serverless caches",
		},
		[]string{"account", "account_name", "region", "engine", "status"},
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
	LambdaFunctions.Reset()
	ECSServices.Reset()
	ECSTasks.Reset()
	ECSServicesByStatus.Reset()
	ECSCapacityProviders.Reset()
	ECSTaskDefinitions.Reset()
	EKSClusters.Reset()
	ELBClassic.Reset()
	ELBV2.Reset()
	ELBV2Detailed.Reset()
	DynamoDBTables.Reset()
	ElastiCacheClusters.Reset()
	ElastiCacheServerlessCaches.Reset()
	OpenSearchDomains.Reset()
	OpenSearchServerlessCollections.Reset()
	MQBrokers.Reset()
	SESIdentities.Reset()
	SESConfigSets.Reset()
	SESContactLists.Reset()
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
}
