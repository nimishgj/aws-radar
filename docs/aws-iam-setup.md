# AWS IAM Setup Guide

This guide explains how to create an IAM user with the minimum required permissions to run AWS Radar.

## Overview

AWS Radar requires **read-only access** to various AWS services to collect resource counts. This guide will help you create a dedicated IAM user with a custom policy that grants only the necessary permissions.

## Prerequisites

- AWS account with administrative access (to create IAM users and policies)
- AWS CLI installed (optional, for CLI-based setup)

## Option 1: AWS Console Setup

### Step 1: Create IAM Policy

1. Sign in to the [AWS Management Console](https://console.aws.amazon.com/)
2. Navigate to **IAM** > **Policies** > **Create policy**
3. Select the **JSON** tab
4. Paste the following policy:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AWSRadarReadOnly",
            "Effect": "Allow",
            "Action": [
                "apigateway:GET",
                "autoscaling:DescribeAutoScalingGroups",
                "athena:ListWorkGroups",
                "apprunner:ListServices",
                "ecr:DescribeRepositories",
                "ec2:DescribeInstances",
                "ec2:DescribeVpcs",
                "ec2:DescribeSubnets",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeNatGateways",
                "ec2:DescribeInternetGateways",
                "ec2:DescribeVolumes",
                "elasticfilesystem:DescribeFileSystems",
                "events:ListEventBuses",
                "events:ListRules",
                "glue:GetJobs",
                "codebuild:ListProjects",
                "codepipeline:ListPipelines",
                "codedeploy:ListApplications",
                "transfer:ListServers",
                "kafka:ListClustersV2",
                "redshift:DescribeClusters",
                "s3:ListAllMyBuckets",
                "s3:GetBucketLocation",
                "rds:DescribeDBInstances",
                "lambda:ListFunctions",
                "ecs:ListClusters",
                "ecs:DescribeClusters",
                "ecs:ListServices",
                "ecs:DescribeServices",
                "ecs:ListTasks",
                "ecs:DescribeTasks",
                "eks:ListClusters",
                "eks:DescribeCluster",
                "elasticloadbalancing:DescribeLoadBalancers",
                "elasticloadbalancingv2:DescribeLoadBalancers",
                "dynamodb:ListTables",
                "dynamodb:DescribeTable",
                "elasticache:DescribeCacheClusters",
                "es:ListDomainNames",
                "guardduty:ListDetectors",
                "securityhub:GetEnabledStandards",
                "inspector2:ListCoverage",
                "macie2:ListClassificationJobs",
                "wafv2:ListWebACLs",
                "secretsmanager:ListSecrets",
                "ssm:DescribeParameters",
                "states:ListStateMachines",
                "sqs:ListQueues",
                "sns:ListTopics",
                "cloudfront:ListDistributions",
                "route53:ListHostedZones",
                "acm:ListCertificates",
                "acm:DescribeCertificate",
                "iam:ListUsers",
                "iam:ListRoles",
                "iam:ListAccountAliases",
                "shield:DescribeSubscription",
                "sts:GetCallerIdentity"
            ],
            "Resource": "*"
        }
    ]
}
```

5. Click **Next**
6. Name the policy: `AWSRadarReadOnlyPolicy`
7. Add description: `Read-only access for AWS Radar monitoring agent`
8. Click **Create policy**

### Step 2: Create IAM User

1. Navigate to **IAM** > **Users** > **Create user**
2. Enter username: `aws-radar`
3. Click **Next**
4. Select **Attach policies directly**
5. Search for and select `AWSRadarReadOnlyPolicy`
6. If you enable `cost_explorer` in AWS Radar config, also attach `AWSRadarCostExplorerPolicy` (create below)
7. Click **Next**
8. Review and click **Create user**

### Step 3: Create Access Keys

1. Click on the newly created user `aws-radar`
2. Go to **Security credentials** tab
3. Under **Access keys**, click **Create access key**
4. Select **Application running outside AWS**
5. Click **Next**
6. (Optional) Add description: `AWS Radar monitoring`
7. Click **Create access key**
8. **Important**: Download or copy the access key ID and secret access key. You won't be able to see the secret key again.

## Option 2: AWS CLI Setup

If you prefer using the AWS CLI, follow these steps:

### Step 1: Create the Policy File

Create or update a file named `aws-radar-policy.json` with the policy shown below:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AWSRadarReadOnly",
            "Effect": "Allow",
            "Action": [
                "apigateway:GET",
                "autoscaling:DescribeAutoScalingGroups",
                "athena:ListWorkGroups",
                "apprunner:ListServices",
                "ecr:DescribeRepositories",
                "ec2:DescribeInstances",
                "ec2:DescribeVpcs",
                "ec2:DescribeSubnets",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeNatGateways",
                "ec2:DescribeInternetGateways",
                "ec2:DescribeVolumes",
                "elasticfilesystem:DescribeFileSystems",
                "events:ListEventBuses",
                "events:ListRules",
                "glue:GetJobs",
                "codebuild:ListProjects",
                "codepipeline:ListPipelines",
                "codedeploy:ListApplications",
                "transfer:ListServers",
                "kafka:ListClustersV2",
                "redshift:DescribeClusters",
                "s3:ListAllMyBuckets",
                "s3:GetBucketLocation",
                "rds:DescribeDBInstances",
                "lambda:ListFunctions",
                "ecs:ListClusters",
                "ecs:DescribeClusters",
                "ecs:ListServices",
                "ecs:DescribeServices",
                "ecs:ListTasks",
                "ecs:DescribeTasks",
                "eks:ListClusters",
                "eks:DescribeCluster",
                "elasticloadbalancing:DescribeLoadBalancers",
                "elasticloadbalancingv2:DescribeLoadBalancers",
                "dynamodb:ListTables",
                "dynamodb:DescribeTable",
                "elasticache:DescribeCacheClusters",
                "es:ListDomainNames",
                "guardduty:ListDetectors",
                "securityhub:GetEnabledStandards",
                "inspector2:ListCoverage",
                "macie2:ListClassificationJobs",
                "wafv2:ListWebACLs",
                "secretsmanager:ListSecrets",
                "ssm:DescribeParameters",
                "states:ListStateMachines",
                "sqs:ListQueues",
                "sns:ListTopics",
                "cloudfront:ListDistributions",
                "route53:ListHostedZones",
                "acm:ListCertificates",
                "acm:DescribeCertificate",
                "iam:ListUsers",
                "iam:ListRoles",
                "iam:ListAccountAliases",
                "shield:DescribeSubscription",
                "sts:GetCallerIdentity"
            ],
            "Resource": "*"
        }
    ]
}
```

### Step 2: Create Policy, User, and Access Keys

```bash
# Create the IAM policy
aws iam create-policy \
    --policy-name AWSRadarReadOnlyPolicy \
    --policy-document file://aws-radar-policy.json \
    --description "Read-only access for AWS Radar monitoring agent"

# Create the IAM user
aws iam create-user --user-name aws-radar

# Get your AWS account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Attach the policy to the user
aws iam attach-user-policy \
    --user-name aws-radar \
    --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSRadarReadOnlyPolicy

# Optional: add Cost Explorer access if cost_explorer is enabled
cat > aws-radar-cost-explorer-policy.json << 'EOF'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowCostExplorerRead",
      "Effect": "Allow",
      "Action": [
        "ce:GetCostAndUsage"
      ],
      "Resource": "*"
    }
  ]
}
EOF

aws iam create-policy \
    --policy-name AWSRadarCostExplorerPolicy \
    --policy-document file://aws-radar-cost-explorer-policy.json \
    --description "Allow AWS Radar cost_explorer collector to read Cost Explorer data"

aws iam attach-user-policy \
    --user-name aws-radar \
    --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSRadarCostExplorerPolicy

# Create access keys
aws iam create-access-key --user-name aws-radar
```

Save the `AccessKeyId` and `SecretAccessKey` from the output.

### Updating an Existing `AWSRadarReadOnlyPolicy`

If the policy already exists and is attached to your IAM user, update it by creating a new default policy version:

```bash
# Get your AWS account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Update existing policy in place
aws iam create-policy-version \
    --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSRadarReadOnlyPolicy \
    --policy-document file://aws-radar-policy.json \
    --set-as-default
```

## Option 3: Using AWS Managed Policies (Alternative)

If you prefer using AWS managed policies instead of a custom policy, you can attach the `ReadOnlyAccess` managed policy. However, this grants broader read access than necessary.

```bash
aws iam attach-user-policy \
    --user-name aws-radar \
    --policy-arn arn:aws:iam::aws:policy/ReadOnlyAccess
```

**Note**: The custom policy in Option 1/2 follows the principle of least privilege and is recommended for production use.

## Configure AWS Radar

Once you have the access keys, configure AWS Radar:

### Using Environment Variables

```bash
cp docker/.env.example docker/.env
```

Edit `docker/.env`:

```bash
AWS_ACCESS_KEY_ID=AKIAXXXXXXXXXXXXXXXX
AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Using AWS Credentials File

Alternatively, configure credentials in `~/.aws/credentials`:

```ini
[aws-radar]
aws_access_key_id = AKIAXXXXXXXXXXXXXXXX
aws_secret_access_key = xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

And set the profile:

```bash
export AWS_PROFILE=aws-radar
```

## Required Permissions Reference

The following table lists all API actions required by AWS Radar:

| Service | Action | Purpose |
|---------|--------|---------|
| **EC2** | `DescribeInstances` | Count EC2 instances |
| | `DescribeVpcs` | Count VPCs |
| | `DescribeSubnets` | Count subnets |
| | `DescribeSecurityGroups` | Count security groups |
| | `DescribeNatGateways` | Count NAT gateways |
| | `DescribeInternetGateways` | Count internet gateways |
| | `DescribeVolumes` | Count EBS volumes |
| **S3** | `ListAllMyBuckets` | List S3 buckets |
| | `GetBucketLocation` | Get bucket region |
| **RDS** | `DescribeDBInstances` | Count RDS instances |
| **Lambda** | `ListFunctions` | Count Lambda functions |
| **ECS** | `ListClusters` | List ECS clusters |
| | `DescribeClusters` | Get cluster details |
| | `ListServices` | List services per cluster |
| | `DescribeServices` | Get service details |
| | `ListTasks` | List tasks per cluster |
| | `DescribeTasks` | Get task details |
| **EKS** | `ListClusters` | List EKS clusters |
| | `DescribeCluster` | Get cluster version |
| **ELB** | `DescribeLoadBalancers` | Count Classic ELBs |
| **ELBv2** | `DescribeLoadBalancers` | Count ALB/NLB |
| **DynamoDB** | `ListTables` | List DynamoDB tables |
| | `DescribeTable` | Get table billing mode |
| **ElastiCache** | `DescribeCacheClusters` | Count ElastiCache clusters |
| **SQS** | `ListQueues` | Count SQS queues |
| **SNS** | `ListTopics` | Count SNS topics |
| **App Runner** | `ListServices` | Count App Runner services |
| **CodeBuild** | `ListProjects` | Count CodeBuild projects |
| **CodePipeline** | `ListPipelines` | Count CodePipeline pipelines |
| **CodeDeploy** | `ListApplications` | Count CodeDeploy applications |
| **Transfer Family** | `ListServers` | Count Transfer Family servers |
| **MSK (Kafka)** | `ListClustersV2` | Count MSK clusters |
| **Redshift** | `DescribeClusters` | Count Redshift clusters |
| **GuardDuty** | `ListDetectors` | Count GuardDuty detectors |
| **Security Hub** | `GetEnabledStandards` | Count enabled Security Hub standards |
| **Inspector2** | `ListCoverage` | Count Inspector2 covered resources |
| **Macie** | `ListClassificationJobs` | Count Macie classification jobs |
| **WAFv2** | `ListWebACLs` | Count regional WAFv2 web ACLs |
| **Shield** | `DescribeSubscription` | Detect Shield Advanced subscription status |
| **CloudFront** | `ListDistributions` | Count CloudFront distributions |
| **Route53** | `ListHostedZones` | Count hosted zones |
| **ACM** | `ListCertificates` | List certificates |
| | `DescribeCertificate` | Get certificate details |
| **IAM** | `ListUsers` | Count IAM users |
| | `ListRoles` | Count IAM roles |
| **Cost Explorer** | `GetCostAndUsage` | Fetch cost by AWS service (attach separate AWSRadarCostExplorerPolicy when cost_explorer is enabled) |

## Security Best Practices

1. **Use least privilege**: The custom policy only grants the minimum required permissions
2. **Rotate access keys**: Regularly rotate the access keys (every 90 days recommended)
3. **Never commit credentials**: Keep credentials out of version control
4. **Use IAM roles when possible**: If running on AWS (EC2, ECS, Lambda), use IAM roles instead of access keys
5. **Enable MFA**: Consider requiring MFA for sensitive operations on the AWS account
6. **Monitor usage**: Enable CloudTrail to monitor API usage by the aws-radar user

## Cleanup

To remove the IAM user and policy:

```bash
# Detach the policy from the user
aws iam detach-user-policy \
    --user-name aws-radar \
    --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSRadarReadOnlyPolicy

# Detach optional Cost Explorer policy (if attached)
aws iam detach-user-policy \
    --user-name aws-radar \
    --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSRadarCostExplorerPolicy

# Delete access keys (list them first)
aws iam list-access-keys --user-name aws-radar
aws iam delete-access-key --user-name aws-radar --access-key-id AKIAXXXXXXXXXXXXXXXX

# Delete the user
aws iam delete-user --user-name aws-radar

# Delete the policy
aws iam delete-policy --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSRadarReadOnlyPolicy

# Delete optional Cost Explorer policy
aws iam delete-policy --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSRadarCostExplorerPolicy
```
