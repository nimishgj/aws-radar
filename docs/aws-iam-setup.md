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
                "mq:ListBrokers",
                "ses:ListEmailIdentities",
                "cloudformation:ListStacks",
                "rds:DescribeDBClusters",
                "memorydb:DescribeClusters",
                "timestream:DescribeEndpoints",
                "timestream:ListDatabases",
                "timestream:ListTables",
                "fsx:DescribeFileSystems",
                "backup:ListBackupVaults",
                "kinesis:ListStreams",
                "firehose:ListDeliveryStreams",
                "kinesisanalytics:ListApplications",
                "elasticmapreduce:ListClusters",
                "elasticbeanstalk:DescribeApplications",
                "kms:ListKeys",
                "cloudtrail:DescribeTrails",
                "batch:DescribeJobQueues",
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
                "ecs:DescribeTaskDefinition",
                "eks:ListClusters",
                "eks:DescribeCluster",
                "elasticloadbalancing:DescribeLoadBalancers",
                "elasticloadbalancingv2:DescribeLoadBalancers",
                "elasticloadbalancingv2:DescribeListeners",
                "elasticloadbalancingv2:DescribeTargetGroups",
                "elasticloadbalancingv2:DescribeRules",
                "dynamodb:ListTables",
                "dynamodb:DescribeTable",
                "elasticache:DescribeCacheClusters",
                "elasticache:DescribeReplicationGroups",
                "elasticache:DescribeGlobalReplicationGroups",
                "es:ListDomainNames",
                "guardduty:ListDetectors",
                "securityhub:GetEnabledStandards",
                "inspector2:ListCoverage",
                "macie2:ListClassificationJobs",
                "wafv2:ListWebACLs",
                "secretsmanager:ListSecrets",
                "ssm:DescribeParameters",
                "ssm:ListDocuments",
                "ssm:DescribeMaintenanceWindows",
                "ssm:ListAssociations",
                "ssm:DescribePatchBaselines",
                "states:ListStateMachines",
                "sqs:ListQueues",
                "sqs:GetQueueAttributes",
                "sns:ListTopics",
                "sns:GetTopicAttributes",
                "cognito-idp:ListUserPools",
                "cognito-idp:DescribeUserPool",
                "cognito-identity:ListIdentityPools",
                "network-firewall:ListFirewalls",
                "network-firewall:ListFirewallPolicies",
                "network-firewall:ListRuleGroups",
                "fms:ListPolicies",
                "acm-pca:ListCertificateAuthorities",
                "servicecatalog:ListPortfolios",
                "servicecatalog:SearchProductsAsAdmin",
                "servicecatalog:SearchProvisionedProducts",
                "license-manager:ListLicenseConfigurations",
                "cloudfront:ListDistributions",
                "route53:ListHostedZones",
                "acm:ListCertificates",
                "acm:DescribeCertificate",
                "iam:ListUsers",
                "iam:ListRoles",
                "iam:ListAccountAliases",
                "shield:DescribeSubscription",
                "autoscaling:DescribePolicies",
                "autoscaling:DescribeLifecycleHooks",
                "autoscaling:DescribeWarmPool",
                "autoscaling:DescribeInstanceRefreshes",
                "ecs:DescribeCapacityProviders",
                "ecs:ListTaskDefinitions",
                "ecr-public:DescribeRepositories",
                "route53resolver:ListResolverEndpoints",
                "route53resolver:ListResolverRules",
                "ec2:DescribeVpcEndpoints",
                "ec2:DescribeTransitGateways",
                "ec2:DescribeVpnGateways",
                "directconnect:DescribeConnections",
                "organizations:ListAccounts",
                "organizations:ListRoots",
                "organizations:ListOrganizationalUnitsForParent",
                "controltower:ListLandingZones",
                "controltower:GetLandingZone",
                "config:DescribeConfigurationRecorders",
                "cloudtrail:ListEventDataStores",
                "kinesisvideo:ListStreams",
                "aoss:ListCollections",
                "ses:ListConfigurationSets",
                "ses:ListContactLists",
                "ses:GetEmailIdentity",
                "ses:GetConfigurationSet",
                "ses:GetConfigurationSetEventDestinations",
                "ses:ListSuppressedDestinations",
                "ses:ListDedicatedIpPools",
                "ses:GetAccount",
                "s3:ListAccessPoints",
                "s3:ListStorageLensConfigurations",
                "rds:DescribeDBProxies",
                "rds:DescribeDBProxyTargets",
                "elasticache:DescribeServerlessCaches",
                "bedrock:ListCustomModels",
                "sagemaker:ListEndpoints",
                "quicksight:ListDashboards",
                "workspaces:DescribeWorkspaces",
                "appstream:DescribeFleets",
                "connect:ListInstances",
                "amplify:ListApps",
                "globalaccelerator:ListAccelerators",
                "datasync:ListTasks",
                "dms:DescribeReplicationInstances",
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
                "mq:ListBrokers",
                "ses:ListEmailIdentities",
                "cloudformation:ListStacks",
                "rds:DescribeDBClusters",
                "memorydb:DescribeClusters",
                "timestream:DescribeEndpoints",
                "timestream:ListDatabases",
                "timestream:ListTables",
                "fsx:DescribeFileSystems",
                "backup:ListBackupVaults",
                "kinesis:ListStreams",
                "firehose:ListDeliveryStreams",
                "kinesisanalytics:ListApplications",
                "elasticmapreduce:ListClusters",
                "elasticbeanstalk:DescribeApplications",
                "kms:ListKeys",
                "cloudtrail:DescribeTrails",
                "batch:DescribeJobQueues",
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
                "ecs:DescribeTaskDefinition",
                "eks:ListClusters",
                "eks:DescribeCluster",
                "elasticloadbalancing:DescribeLoadBalancers",
                "elasticloadbalancingv2:DescribeLoadBalancers",
                "elasticloadbalancingv2:DescribeListeners",
                "elasticloadbalancingv2:DescribeTargetGroups",
                "elasticloadbalancingv2:DescribeRules",
                "dynamodb:ListTables",
                "dynamodb:DescribeTable",
                "elasticache:DescribeCacheClusters",
                "elasticache:DescribeReplicationGroups",
                "elasticache:DescribeGlobalReplicationGroups",
                "es:ListDomainNames",
                "guardduty:ListDetectors",
                "securityhub:GetEnabledStandards",
                "inspector2:ListCoverage",
                "macie2:ListClassificationJobs",
                "wafv2:ListWebACLs",
                "secretsmanager:ListSecrets",
                "ssm:DescribeParameters",
                "ssm:ListDocuments",
                "ssm:DescribeMaintenanceWindows",
                "ssm:ListAssociations",
                "ssm:DescribePatchBaselines",
                "states:ListStateMachines",
                "sqs:ListQueues",
                "sqs:GetQueueAttributes",
                "sns:ListTopics",
                "sns:GetTopicAttributes",
                "cognito-idp:ListUserPools",
                "cognito-idp:DescribeUserPool",
                "cognito-identity:ListIdentityPools",
                "network-firewall:ListFirewalls",
                "network-firewall:ListFirewallPolicies",
                "network-firewall:ListRuleGroups",
                "fms:ListPolicies",
                "acm-pca:ListCertificateAuthorities",
                "servicecatalog:ListPortfolios",
                "servicecatalog:SearchProductsAsAdmin",
                "servicecatalog:SearchProvisionedProducts",
                "license-manager:ListLicenseConfigurations",
                "cloudfront:ListDistributions",
                "route53:ListHostedZones",
                "acm:ListCertificates",
                "acm:DescribeCertificate",
                "iam:ListUsers",
                "iam:ListRoles",
                "iam:ListAccountAliases",
                "shield:DescribeSubscription",
                "autoscaling:DescribePolicies",
                "autoscaling:DescribeLifecycleHooks",
                "autoscaling:DescribeWarmPool",
                "autoscaling:DescribeInstanceRefreshes",
                "ecs:DescribeCapacityProviders",
                "ecs:ListTaskDefinitions",
                "ecr-public:DescribeRepositories",
                "route53resolver:ListResolverEndpoints",
                "route53resolver:ListResolverRules",
                "ec2:DescribeVpcEndpoints",
                "ec2:DescribeTransitGateways",
                "ec2:DescribeVpnGateways",
                "directconnect:DescribeConnections",
                "organizations:ListAccounts",
                "organizations:ListRoots",
                "organizations:ListOrganizationalUnitsForParent",
                "controltower:ListLandingZones",
                "controltower:GetLandingZone",
                "config:DescribeConfigurationRecorders",
                "cloudtrail:ListEventDataStores",
                "kinesisvideo:ListStreams",
                "aoss:ListCollections",
                "ses:ListConfigurationSets",
                "ses:ListContactLists",
                "ses:GetEmailIdentity",
                "ses:GetConfigurationSet",
                "ses:GetConfigurationSetEventDestinations",
                "ses:ListSuppressedDestinations",
                "ses:ListDedicatedIpPools",
                "ses:GetAccount",
                "s3:ListAccessPoints",
                "s3:ListStorageLensConfigurations",
                "rds:DescribeDBProxies",
                "rds:DescribeDBProxyTargets",
                "elasticache:DescribeServerlessCaches",
                "bedrock:ListCustomModels",
                "sagemaker:ListEndpoints",
                "quicksight:ListDashboards",
                "workspaces:DescribeWorkspaces",
                "appstream:DescribeFleets",
                "connect:ListInstances",
                "amplify:ListApps",
                "globalaccelerator:ListAccelerators",
                "datasync:ListTasks",
                "dms:DescribeReplicationInstances",
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
| | `DescribeListeners` | Count listeners by protocol |
| | `DescribeTargetGroups` | Count target groups by target type |
| | `DescribeRules` | Count listener rules per ALB |
| **DynamoDB** | `ListTables` | List DynamoDB tables |
| | `DescribeTable` | Get table billing mode |
| **ElastiCache** | `DescribeCacheClusters` | Count ElastiCache clusters |
| **SQS** | `ListQueues` | Count SQS queues |
| | `GetQueueAttributes` | Get message counts, DLQ config per queue |
| **SNS** | `ListTopics` | Count SNS topics |
| | `GetTopicAttributes` | Get subscription counts, FIFO status per topic |
| **MQ** | `ListBrokers` | Count Amazon MQ brokers |
| **SES** | `ListEmailIdentities` | Count SES identities |
| **CloudFormation** | `ListStacks` | Count CloudFormation stacks |
| **DocumentDB / Neptune** | `DescribeDBClusters` | Count DocumentDB and Neptune clusters |
| **MemoryDB** | `DescribeClusters` | Count MemoryDB clusters |
| **Timestream** | `DescribeEndpoints` | Resolve endpoint for Timestream APIs |
| | `ListDatabases` | Count Timestream databases |
| | `ListTables` | Count Timestream tables |
| **FSx** | `DescribeFileSystems` | Count FSx file systems |
| **Backup** | `ListBackupVaults` | Count AWS Backup vaults |
| **Kinesis Data Streams** | `ListStreams` | Count Kinesis streams |
| **Kinesis Data Firehose** | `ListDeliveryStreams` | Count Firehose delivery streams |
| **Kinesis Data Analytics** | `ListApplications` | Count Kinesis Analytics applications |
| **EMR** | `ListClusters` | Count EMR clusters |
| **Elastic Beanstalk** | `DescribeApplications` | Count Elastic Beanstalk applications |
| **KMS** | `ListKeys` | Count KMS keys |
| **CloudTrail** | `DescribeTrails` | Count CloudTrail trails |
| **AWS Batch** | `DescribeJobQueues` | Count AWS Batch job queues |
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
| **Route53 Resolver** | `ListResolverEndpoints` | Count Route53 Resolver endpoints |
| | `ListResolverRules` | Count Route53 Resolver rules |
| **ACM** | `ListCertificates` | List certificates |
| | `DescribeCertificate` | Get certificate details |
| **IAM** | `ListUsers` | Count IAM users |
| | `ListRoles` | Count IAM roles |
| **Auto Scaling** | `DescribePolicies` | Count Auto Scaling policies |
| | `DescribeLifecycleHooks` | Count lifecycle hooks per Auto Scaling group |
| | `DescribeWarmPool` | Count warm pool groups and instances |
| | `DescribeInstanceRefreshes` | Count instance refreshes by status |
| **ECS** | `DescribeCapacityProviders` | Count ECS capacity providers |
| | `ListTaskDefinitions` | Count ECS task definitions |
| | `DescribeTaskDefinition` | Classify task definitions by family/revision/runtime platform |
| **ECR Public** | `DescribeRepositories` | Count ECR Public repositories |
| **VPC / EC2 networking** | `DescribeVpcEndpoints` | Count VPC endpoints |
| | `DescribeTransitGateways` | Count transit gateways |
| | `DescribeVpnGateways` | Count VPN gateways |
| **Direct Connect** | `DescribeConnections` | Count Direct Connect connections |
| **Organizations** | `ListAccounts` | Count organization accounts by state |
| | `ListRoots` | Discover organization roots |
| | `ListOrganizationalUnitsForParent` | Count organizational units |
| **Control Tower** | `ListLandingZones` | List landing zones |
| | `GetLandingZone` | Read landing zone deployment status |
| **Config** | `DescribeConfigurationRecorders` | Count configuration recorders |
| **CloudTrail Lake** | `ListEventDataStores` | Count CloudTrail Lake event data stores |
| **Kinesis Video Streams** | `ListStreams` | Count Kinesis Video streams |
| **OpenSearch Serverless** | `ListCollections` | Count serverless collections |
| **SES** | `ListConfigurationSets` | Count SES config sets |
| | `ListContactLists` | Count SES contact lists |
| | `GetEmailIdentity` | Classify identities by verification/DKIM/mail-from status |
| | `GetConfigurationSet` | Read dedicated IP pool association per configuration set |
| | `GetConfigurationSetEventDestinations` | Count event destinations by type |
| | `ListSuppressedDestinations` | Count suppression list entries by reason |
| | `ListDedicatedIpPools` | Count dedicated IP pools |
| | `GetAccount` | Read sending enabled/production access/quota values |
| **S3 (control plane)** | `ListAccessPoints` | Count S3 access points |
| | `ListStorageLensConfigurations` | Count S3 Storage Lens configs |
| **RDS** | `DescribeDBProxies` | Count RDS proxies |
| | `DescribeDBProxyTargets` | Count RDS proxy targets by type |
| **ElastiCache** | `DescribeServerlessCaches` | Count ElastiCache serverless caches |
| | `DescribeReplicationGroups` | Count replication groups by engine/status/cluster mode |
| | `DescribeGlobalReplicationGroups` | Count global replication groups by status |
| **Bedrock** | `ListCustomModels` | Count Bedrock custom models |
| **SageMaker** | `ListEndpoints` | Count SageMaker endpoints |
| **QuickSight** | `ListDashboards` | Count QuickSight dashboards |
| **WorkSpaces** | `DescribeWorkspaces` | Count WorkSpaces by state |
| **AppStream 2.0** | `DescribeFleets` | Count AppStream fleets by state |
| **Connect** | `ListInstances` | Count Connect instances |
| **Amplify** | `ListApps` | Count Amplify apps |
| **Global Accelerator** | `ListAccelerators` | Count global accelerators |
| **DataSync** | `ListTasks` | Count DataSync tasks |
| **DMS** | `DescribeReplicationInstances` | Count DMS replication instances |
| **SSM** | `DescribeParameters` | Count SSM parameters by type |
| | `ListDocuments` | Count SSM documents by owner |
| | `DescribeMaintenanceWindows` | Count maintenance windows by state |
| | `ListAssociations` | Count SSM associations |
| | `DescribePatchBaselines` | Count patch baselines |
| **Cognito** | `ListUserPools` | List Cognito user pools |
| | `DescribeUserPool` | Get estimated user count per pool |
| | `ListIdentityPools` | Count Cognito identity pools |
| **Network Firewall** | `ListFirewalls` | Count Network Firewall firewalls |
| | `ListFirewallPolicies` | Count firewall policies |
| | `ListRuleGroups` | Count rule groups |
| **Firewall Manager** | `ListPolicies` | Count FMS policies by resource type |
| **ACM PCA** | `ListCertificateAuthorities` | Count private certificate authorities by status/type |
| **Service Catalog** | `ListPortfolios` | Count Service Catalog portfolios |
| | `SearchProductsAsAdmin` | Count Service Catalog products |
| | `SearchProvisionedProducts` | Count provisioned products by status |
| **License Manager** | `ListLicenseConfigurations` | Count license configurations |
| **Cost Explorer** | `GetCostAndUsage` | Fetch cost by AWS service (attach separate AWSRadarCostExplorerPolicy when cost_explorer is enabled) |

## Security Best Practices

1. **Use least privilege**: The custom policy only grants the minimum required permissions
2. **Rotate access keys**: Regularly rotate the access keys (every 90 days recommended)
3. **Never commit credentials**: Keep credentials out of version control
4. **Use IAM roles when possible**: If running on AWS (EC2, ECS, Lambda), use IAM roles instead of access keys
5. **Enable MFA**: Consider requiring MFA for sensitive operations on the AWS account
6. **Monitor usage**: Enable CloudTrail to monitor API usage by the aws-radar user

## Advanced Dimension Metrics (ELB/ECS/ASG/RDS/SES)

When the latest collectors are enabled, AWS Radar also emits these additional metrics used by the advanced Grafana panels:

- `aws_elbv2_listeners_total`
- `aws_elbv2_target_groups_total`
- `aws_elbv2_rules_per_alb`
- `aws_elbv2_availability_zones_per_lb`
- `aws_elbv2_subnets_per_lb`
- `aws_ecs_tasks_by_status_total`
- `aws_ecs_cluster_depth`
- `aws_ecs_capacity_providers_detailed_total`
- `aws_ecs_default_capacity_provider_strategy_total`
- `aws_ecs_task_definitions_detailed_total`
- `aws_autoscaling_policies_by_type_total`
- `aws_autoscaling_groups_by_mixed_instances_total`
- `aws_autoscaling_launch_template_usage_total`
- `aws_autoscaling_lifecycle_hooks_total`
- `aws_autoscaling_warm_pools_total`
- `aws_autoscaling_warm_pool_instances_total`
- `aws_autoscaling_instance_refreshes_total`
- `aws_rds_instances_by_engine_version_total`
- `aws_rds_instances_by_class_total`
- `aws_rds_instances_by_multi_az_total`
- `aws_rds_read_replicas_total`
- `aws_rds_proxy_targets_total`
- `aws_rds_aurora_serverless_v2_capacity`
- `aws_rds_aurora_serverless_by_status_total`
- `aws_elasticache_replication_groups_total`
- `aws_elasticache_global_replication_groups_total`
- `aws_ses_identities_by_verification_status_total`
- `aws_ses_identity_auth_status_total`
- `aws_ses_config_set_event_destinations_total`
- `aws_ses_suppressed_destinations_total`
- `aws_ses_dedicated_ip_pools_total`
- `aws_ses_config_sets_by_sending_pool_total`
- `aws_ses_account_settings`
- `aws_ses_sending_quota`
- `aws_sns_subscriptions_total`
- `aws_sns_topics_by_type_total`
- `aws_sqs_messages_total`
- `aws_sqs_queues_with_dlq_total`
- `aws_ssm_documents_total`
- `aws_ssm_maintenance_windows_total`
- `aws_ssm_associations_total`
- `aws_ssm_patch_baselines_total`
- `aws_cognito_user_pools_total`
- `aws_cognito_user_pool_users_total`
- `aws_cognito_identity_pools_total`
- `aws_networkfirewall_firewalls_total`
- `aws_networkfirewall_policies_total`
- `aws_networkfirewall_rule_groups_total`
- `aws_fms_policies_total`
- `aws_acmpca_certificate_authorities_total`
- `aws_servicecatalog_portfolios_total`
- `aws_servicecatalog_products_total`
- `aws_servicecatalog_provisioned_products_total`
- `aws_licensemanager_license_configurations_total`

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
