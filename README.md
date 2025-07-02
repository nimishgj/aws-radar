# AWS Radar

## Cloud Resource Inventory Tool for AWS

AWS Radar is a comprehensive resource inventory tool that automatically discovers and catalogs AWS resources across all regions. It collects information about various AWS resources and outputs them to a single CSV file, making it easy to audit, document, and track your AWS infrastructure.

## Features

- **Comprehensive Coverage**: Inventories 70+ types of AWS resources across multiple service categories
- **Multi-Region Support**: Scans all available AWS regions
- **Read-Only Operations**: Uses only AWS CLI list and describe commands (no modifications)
- **Cost Conscious**: Avoids operations that could incur AWS charges
- **Modular Architecture**: Organized by service category for easy maintenance and extension
- **Dynamic Script Discovery**: Automatically finds and executes all inventory scripts
- **CSV Output**: Produces a standardized output for easy analysis

## Prerequisites

### AWS CLI

- AWS CLI version 2.x installed and configured
- To install: Follow the [official AWS CLI installation guide](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

### AWS IAM Permissions

For read-only inventory collection, your AWS credentials should have the following permissions:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:Describe*",
                "rds:Describe*",
                "s3:List*",
                "s3:GetBucket*",
                "lambda:List*",
                "lambda:Get*",
                "ecs:List*",
                "ecs:Describe*",
                "ecr:Describe*",
                "elasticache:Describe*",
                "cloudwatch:Describe*",
                "cloudwatch:List*",
                "cloudwatch:Get*",
                "logs:Describe*",
                "logs:List*",
                "dynamodb:List*",
                "dynamodb:Describe*",
                "sns:List*",
                "sqs:List*",
                "mq:List*",
                "mq:Describe*",
                "kafka:List*",
                "secretsmanager:List*",
                "kms:List*",
                "ssm:Describe*",
                "route53:List*",
                "apigateway:GET",
                "eks:List*",
                "eks:Describe*",
                "kinesis:List*",
                "kinesis:Describe*",
                "states:List*"
            ],
            "Resource": "*"
        }
    ]
}
```

You can attach this policy to an IAM role or user and use those credentials with the AWS CLI.

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/nimishgj/aws-radar.git
   cd aws-radar
   ```

2. Ensure all scripts are executable:
   ```bash
   find . -name "*.sh" -exec chmod +x {} \;
   ```

## Usage

### Basic Usage

Run the main script to collect inventory of all AWS resources:

```bash
./main.sh
```

This will:
1. Create a file named `aws_resources.csv` in the current directory
2. Discover all AWS regions available to your account
3. Scan each region for resources across all service categories
4. Output the results to the CSV file

### Output Format

The output CSV file contains the following columns:
- Resource Type
- Resource ID/Name
- AWS Region

Example:
```
EC2 Instance,i-1234567890abcdef0,us-east-1
S3 Bucket,my-bucket,global
RDS Instance,my-database,us-west-2
```

## How It Works

### Main Script Structure

The `main.sh` script orchestrates the entire inventory process:

1. **Initialization**: Creates the output CSV file and fetches all AWS regions
2. **Dynamic Script Discovery**: Finds all service directories and their scripts
3. **Script Execution**: Runs each script with appropriate parameters
4. **Output Consolidation**: All scripts append to the same CSV file

```bash
# Key parts of main.sh
# Initialize the output file
> "$OUTPUT_FILE"

# Get all AWS regions
AWS_REGIONS=($(aws ec2 describe-regions --query "Regions[].RegionName" --output text))

# Dynamic script execution
for SERVICE_DIR in */; do
  SERVICE_NAME=${SERVICE_DIR%/}
  # Find and execute all scripts in each service directory
  # ...
done
```

### Directory Structure

AWS Radar is organized into service-specific directories:

```
aws-radar/
├── README.md
├── main.sh
├── ec2/
│   ├── instances.sh
│   ├── ebs_volumes.sh
│   ├── ...
├── s3/
│   └── s3.sh
├── vpc/
│   ├── vpc.sh
│   ├── subnets.sh
│   ├── ...
└── ...
```

Each directory contains scripts for collecting specific resource types within that service category.

### Script Pattern

All resource scripts follow a consistent pattern:

```bash
#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Service Resource Name"

# Loop through all regions
for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching $RESOURCE_TYPE in $REGION"
  
  # Use AWS CLI to list resources
  RESOURCES=$(aws service list-resources --region "$REGION" ...)
  
  # Process and format results
  # ...
  
  # Append to output file
  echo "$RESOURCE_TYPE,$RESOURCE_ID,$REGION" >> "$OUTPUT_FILE"
done
```

Special cases:
- Global services like S3 may not need the region parameter
- Some scripts may require additional API calls to gather complete information

## Supported AWS Services

AWS Radar currently supports inventory collection for the following service categories:

1. **EC2**: Instances, Volumes, Snapshots, Security Groups, Elastic IPs, etc.
2. **S3**: Buckets
3. **VPC**: VPCs, Subnets, Route Tables, Gateways, ACLs, etc.
4. **ECS**: Clusters, Task Definitions, Namespaces
5. **ECR**: Public and Private Repositories
6. **RDS**: Instances, Snapshots
7. **CloudWatch**: Log Groups, Dashboards, Alarms, and other monitoring services
8. **Lambda**: Functions, Layers, Event Source Mappings, Function URLs
9. **ElastiCache**: Clusters, Replication Groups, Parameter Groups, etc.
10. **Amazon MQ**: Brokers
11. **MSK**: Clusters, Configurations, Connectors, Replicators
12. **SQS**: Queues
13. **API Gateway**: REST, HTTP, and WebSocket APIs
14. **EKS**: Clusters
15. **SNS**: Topics
16. **Secrets Manager**: Secrets
17. **KMS**: Keys
18. **SSM**: Parameters
19. **Route53**: Hosted Zones
20. **DynamoDB**: Tables
21. **Kinesis**: Streams
22. **Step Functions**: State Machines

## Extending AWS Radar

### Adding a New Service Category

To add inventory for a new AWS service category:

1. Create a new directory for the service:
   ```bash
   mkdir new_service
   ```

2. Create scripts for each resource type:
   ```bash
   touch new_service/resource_type.sh
   chmod +x new_service/resource_type.sh
   ```

3. Follow the standard script pattern (see examples in existing directories)

The main script will automatically discover and execute your new scripts.

### Adding Resources to an Existing Service

Simply add a new script in the appropriate service directory, following the standard pattern.

## Examples

### Running a Full Inventory

```bash
./main.sh
```

### Viewing Results

```bash
cat aws_resources.csv
```

### Filtering Results

```bash
# Show only EC2 resources
grep "EC2" aws_resources.csv

# Show resources in a specific region
grep "us-east-1" aws_resources.csv
```

## Limitations

- AWS Radar performs read-only operations, but excessive API calls may hit AWS rate limits
- Some resource types may require additional API calls, slowing down the inventory process
- Global services like IAM are not fully supported in the current version


