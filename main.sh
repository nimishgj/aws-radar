#!/bin/bash

# AWS Radar - AWS Resource Inventory Tool
# Output file name
OUTPUT_FILE="aws_resources.csv"

# Initialize the output file
> "$OUTPUT_FILE"
echo "INFO: Initialized $OUTPUT_FILE"

# Get all AWS regions
echo "INFO: Fetching AWS regions"
AWS_REGIONS=($(aws ec2 describe-regions --query "Regions[].RegionName" --output text))

# EC2 resources
echo "INFO: Processing EC2 resources"
./ec2/instances.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ec2/ebs_volumes.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ec2/ebs_snapshots.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ec2/security_group.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ec2/elastic_ips.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ec2/key_pairs.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ec2/network_interfaces.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ec2/load_balancers.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ec2/target_groups.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ec2/asg.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# S3 resources
echo "INFO: Processing S3 resources"
./s3/s3.sh "$OUTPUT_FILE"

# VPC resources
echo "INFO: Processing VPC resources"
./vpc/vpc.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/subnets.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/route_table.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/egress_igw.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/dhpc_option_set.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/managed_prefix_list.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/nat_gateway.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/network_acls.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/customer_gateway.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/virtual_private_gateway.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/site_to_site_vpn.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./vpc/transit_gateways.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# ECS resources
echo "INFO: Processing ECS resources"
./ecs/cluster.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ecs/namespaces.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ecs/task_definitions.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# ECR resources
echo "INFO: Processing ECR resources"
./ecr/public_repositories.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./ecr/private_repositories.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# RDS resources
echo "INFO: Processing RDS resources"
./rds/rds.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./rds/snapshots.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# CloudWatch resources
echo "INFO: Processing CloudWatch resources"
./cloudwatch/log_groups.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/dashboards.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/alarms.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/anomaly_detectors.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/synthetics_canaries.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/contributor_insights.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/evidently.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/rum.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/servicelens.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/internet_monitor.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./cloudwatch/logs_insights.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Lambda resources
echo "INFO: Processing Lambda resources"
./lambda/functions.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./lambda/layers.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./lambda/event_source_mappings.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./lambda/function_urls.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# ElastiCache resources
echo "INFO: Processing ElastiCache resources"
./elasticache/clusters.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./elasticache/replication_groups.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./elasticache/parameter_groups.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./elasticache/subnet_groups.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./elasticache/security_groups.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./elasticache/valkey.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./elasticache/memcached.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./elasticache/redis_oss.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./elasticache/global_datastores.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Amazon MQ resources
echo "INFO: Processing Amazon MQ resources"
./amazon_mq/brokers.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# MSK resources
echo "INFO: Processing MSK resources"
./msk/clusters.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./msk/configurations.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./msk/connectors.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./msk/replicators.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# SQS resources
echo "INFO: Processing SQS resources"
./sqs/queues.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# API Gateway resources
echo "INFO: Processing API Gateway resources"
./apigateway/rest_apis.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./apigateway/http_apis.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
./apigateway/websocket_apis.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# EKS resources
echo "INFO: Processing EKS resources"
./eks/clusters.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# SNS resources
echo "INFO: Processing SNS resources"
./sns/topics.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Secrets Manager resources
echo "INFO: Processing Secrets Manager resources"
./secretsmanager/secrets.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# KMS resources
echo "INFO: Processing KMS resources"
./kms/keys.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# SSM resources
echo "INFO: Processing SSM resources"
./ssm/parameters.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Route53 resources
echo "INFO: Processing Route53 resources"
./route53/hosted_zones.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# DynamoDB resources
echo "INFO: Processing DynamoDB resources"
./dynamodb/tables.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Kinesis resources
echo "INFO: Processing Kinesis resources"
./kinesis/streams.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Step Functions resources
echo "INFO: Processing Step Functions resources"
./stepfunctions/state_machines.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

echo "INFO: AWS Radar inventory completed successfully."
echo "INFO: Results saved to $OUTPUT_FILE"
