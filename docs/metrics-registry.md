# Metrics Registry

This document describes all metrics emitted by AWS Radar, including their attributes (labels) and types.

## Overview

All AWS resource metrics:
- Are prefixed with `aws_`
- Are of type **Gauge** (representing current resource counts)
- Are reset before each collection cycle to ensure accuracy

## AWS Resource Metrics

### Compute Services

#### EC2 Instances

| Metric | `aws_ec2_instances_total` |
|--------|---------------------------|
| **Type** | Gauge |
| **Description** | Total number of EC2 instances |
| **Labels** | |
| `region` | AWS region (e.g., `us-east-1`) |
| `instance_type` | EC2 instance type (e.g., `t3.micro`, `m5.large`) |
| `state` | Instance state (`running`, `stopped`, `pending`, etc.) |
| `availability_zone` | Availability zone (e.g., `us-east-1a`) |

#### Lambda Functions

| Metric | `aws_lambda_functions_total` |
|--------|------------------------------|
| **Type** | Gauge |
| **Description** | Total number of Lambda functions |
| **Labels** | |
| `region` | AWS region |
| `runtime` | Function runtime (e.g., `python3.9`, `nodejs18.x`) |
| `memory_size` | Configured memory in MB |

### Container Services

#### ECS Services

| Metric | `aws_ecs_services_total` |
|--------|--------------------------|
| **Type** | Gauge |
| **Description** | Total number of ECS services |
| **Labels** | |
| `region` | AWS region |
| `cluster_name` | ECS cluster name |
| `launch_type` | Launch type (`EC2`, `FARGATE`) |

#### ECS Tasks

| Metric | `aws_ecs_tasks_total` |
|--------|----------------------|
| **Type** | Gauge |
| **Description** | Total number of ECS tasks |
| **Labels** | |
| `region` | AWS region |
| `cluster_name` | ECS cluster name |
| `launch_type` | Launch type (`EC2`, `FARGATE`) |

#### EKS Clusters

| Metric | `aws_eks_clusters_total` |
|--------|--------------------------|
| **Type** | Gauge |
| **Description** | Total number of EKS clusters |
| **Labels** | |
| `region` | AWS region |
| `version` | Kubernetes version |

### Database Services

#### RDS Instances

| Metric | `aws_rds_instances_total` |
|--------|---------------------------|
| **Type** | Gauge |
| **Description** | Total number of RDS instances |
| **Labels** | |
| `region` | AWS region |
| `db_instance_class` | Instance class (e.g., `db.t3.micro`) |
| `engine` | Database engine (`mysql`, `postgres`, `aurora`, etc.) |
| `multi_az` | Multi-AZ deployment (`true`, `false`) |

#### DynamoDB Tables

| Metric | `aws_dynamodb_tables_total` |
|--------|----------------------------|
| **Type** | Gauge |
| **Description** | Total number of DynamoDB tables |
| **Labels** | |
| `region` | AWS region |
| `billing_mode` | Billing mode (`PROVISIONED`, `PAY_PER_REQUEST`) |

#### ElastiCache Clusters

| Metric | `aws_elasticache_clusters_total` |
|--------|----------------------------------|
| **Type** | Gauge |
| **Description** | Total number of ElastiCache clusters |
| **Labels** | |
| `region` | AWS region |
| `engine` | Cache engine (`redis`, `memcached`) |
| `cache_node_type` | Node type (e.g., `cache.t3.micro`) |

### Storage Services

#### S3 Buckets

| Metric | `aws_s3_buckets_total` |
|--------|------------------------|
| **Type** | Gauge |
| **Description** | Total number of S3 buckets |
| **Labels** | |
| `region` | AWS region where bucket was created |

#### EBS Volumes

| Metric | `aws_ebs_volumes_total` |
|--------|-------------------------|
| **Type** | Gauge |
| **Description** | Total number of EBS volumes |
| **Labels** | |
| `region` | AWS region |
| `volume_type` | Volume type (`gp2`, `gp3`, `io1`, `io2`, `st1`, `sc1`) |
| `state` | Volume state (`available`, `in-use`, `creating`, etc.) |

### Networking Services

#### VPCs

| Metric | `aws_vpc_total` |
|--------|-----------------|
| **Type** | Gauge |
| **Description** | Total number of VPCs |
| **Labels** | |
| `region` | AWS region |
| `state` | VPC state (`available`, `pending`) |

#### Subnets

| Metric | `aws_subnet_total` |
|--------|--------------------|
| **Type** | Gauge |
| **Description** | Total number of subnets |
| **Labels** | |
| `region` | AWS region |
| `availability_zone` | Availability zone |

#### Security Groups

| Metric | `aws_security_groups_total` |
|--------|----------------------------|
| **Type** | Gauge |
| **Description** | Total number of security groups |
| **Labels** | |
| `region` | AWS region |
| `vpc_id` | VPC ID the security group belongs to |

#### NAT Gateways

| Metric | `aws_nat_gateways_total` |
|--------|--------------------------|
| **Type** | Gauge |
| **Description** | Total number of NAT gateways |
| **Labels** | |
| `region` | AWS region |
| `state` | Gateway state (`available`, `pending`, `deleting`, etc.) |

#### Internet Gateways

| Metric | `aws_internet_gateways_total` |
|--------|-------------------------------|
| **Type** | Gauge |
| **Description** | Total number of Internet gateways |
| **Labels** | |
| `region` | AWS region |

### Load Balancing

#### Classic Load Balancers

| Metric | `aws_elb_classic_total` |
|--------|-------------------------|
| **Type** | Gauge |
| **Description** | Total number of Classic Load Balancers |
| **Labels** | |
| `region` | AWS region |
| `scheme` | Load balancer scheme (`internet-facing`, `internal`) |

#### Application/Network Load Balancers

| Metric | `aws_elbv2_total` |
|--------|-------------------|
| **Type** | Gauge |
| **Description** | Total number of ALB/NLB load balancers |
| **Labels** | |
| `region` | AWS region |
| `type` | Load balancer type (`application`, `network`, `gateway`) |
| `scheme` | Load balancer scheme (`internet-facing`, `internal`) |

### Messaging Services

#### SQS Queues

| Metric | `aws_sqs_queues_total` |
|--------|------------------------|
| **Type** | Gauge |
| **Description** | Total number of SQS queues |
| **Labels** | |
| `region` | AWS region |
| `queue_type` | Queue type (`standard`, `fifo`) |

#### SNS Topics

| Metric | `aws_sns_topics_total` |
|--------|------------------------|
| **Type** | Gauge |
| **Description** | Total number of SNS topics |
| **Labels** | |
| `region` | AWS region |

### Content Delivery

#### CloudFront Distributions

| Metric | `aws_cloudfront_distributions_total` |
|--------|--------------------------------------|
| **Type** | Gauge |
| **Description** | Total number of CloudFront distributions |
| **Labels** | |
| `price_class` | Price class (`PriceClass_All`, `PriceClass_200`, `PriceClass_100`) |
| `enabled` | Distribution enabled (`true`, `false`) |

### DNS Services

#### Route53 Hosted Zones

| Metric | `aws_route53_hosted_zones_total` |
|--------|----------------------------------|
| **Type** | Gauge |
| **Description** | Total number of Route53 hosted zones |
| **Labels** | None (global service) |

### Security Services

#### ACM Certificates

| Metric | `aws_acm_certificates_total` |
|--------|------------------------------|
| **Type** | Gauge |
| **Description** | Total number of ACM certificates |
| **Labels** | |
| `region` | AWS region |
| `status` | Certificate status (`ISSUED`, `PENDING_VALIDATION`, `EXPIRED`, etc.) |
| `type` | Certificate type (`AMAZON_ISSUED`, `IMPORTED`) |

### Identity Services

#### IAM Users

| Metric | `aws_iam_users_total` |
|--------|----------------------|
| **Type** | Gauge |
| **Description** | Total number of IAM users |
| **Labels** | None (global service) |

#### IAM Roles

| Metric | `aws_iam_roles_total` |
|--------|----------------------|
| **Type** | Gauge |
| **Description** | Total number of IAM roles |
| **Labels** | None (global service) |

## Internal Metrics

These metrics provide observability into AWS Radar's own operation.

### Collection Duration

| Metric | `aws_radar_collection_duration_seconds` |
|--------|----------------------------------------|
| **Type** | Histogram |
| **Description** | Duration of AWS resource collection |
| **Labels** | |
| `collector` | Name of the collector (e.g., `ec2`, `s3`, `rds`) |
| **Buckets** | Default Prometheus buckets (0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10) |

### Collection Errors

| Metric | `aws_radar_collection_errors_total` |
|--------|-------------------------------------|
| **Type** | Counter |
| **Description** | Total number of collection errors |
| **Labels** | |
| `collector` | Name of the collector |
| `region` | AWS region where the error occurred |

## Example Queries

### PromQL Examples

```promql
# Total EC2 instances across all regions
sum(aws_ec2_instances_total)

# EC2 instances by region
sum by (region) (aws_ec2_instances_total)

# Running EC2 instances only
sum(aws_ec2_instances_total{state="running"})

# Lambda functions by runtime
sum by (runtime) (aws_lambda_functions_total)

# Collection error rate by collector
rate(aws_radar_collection_errors_total[5m])
```

### ClickHouse Examples (via OTLP)

When metrics are stored in ClickHouse via OpenTelemetry:

```sql
-- Total EC2 instances per region
SELECT
    Attributes['region'] as region,
    sum(Value) as total
FROM otel_metrics_gauge
WHERE MetricName = 'aws_ec2_instances_total'
    AND TimeUnix = (SELECT max(TimeUnix) FROM otel_metrics_gauge WHERE MetricName = 'aws_ec2_instances_total')
GROUP BY region;

-- EC2 instances over time
SELECT
    toStartOfMinute(fromUnixTimestamp64Milli(TimeUnix)) as time,
    Attributes['region'] as region,
    sum(Value) as total
FROM otel_metrics_gauge
WHERE MetricName = 'aws_ec2_instances_total'
    AND TimeUnix >= toUnixTimestamp64Milli(now() - INTERVAL 1 HOUR)
GROUP BY time, region
ORDER BY time;
```
