# AWS Radar

A lightweight AWS resource monitoring agent that collects resource counts across all AWS regions and exposes them as Prometheus metrics.

## Overview

AWS Radar periodically scans your AWS account and collects resource counts for various services including EC2, Lambda, S3, RDS, VPC, IAM, and more. The metrics are exposed via a Prometheus endpoint, which can be scraped by OpenTelemetry Collector and stored in ClickHouse for visualization in Grafana.

## Architecture

```
┌─────────────┐     ┌──────────────────┐     ┌────────────┐     ┌─────────┐
│  AWS APIs   │────▶│   AWS Radar      │────▶│   OTel     │────▶│ ClickHouse│
│             │     │  (Prometheus     │     │ Collector  │     │         │
│             │     │   :9090/metrics) │     │            │     │         │
└─────────────┘     └──────────────────┘     └────────────┘     └────┬────┘
                                                                      │
                                                                      ▼
                                                                ┌─────────┐
                                                                │ Grafana │
                                                                │  :3000  │
                                                                └─────────┘
```

## Features

- Collects resource counts from 16+ AWS services
- Supports all AWS regions (configurable)
- Exposes Prometheus metrics endpoint
- Pre-built Grafana dashboard with time series visualization
- Docker Compose setup for easy deployment
- Lightweight and efficient

## Supported AWS Services

| Service | Metrics |
|---------|---------|
| EC2 | Instance count by region, type, state |
| Lambda | Function count by region, runtime |
| S3 | Bucket count by region |
| RDS | Instance count by region, engine, class |
| DynamoDB | Table count by region, billing mode |
| EBS | Volume count by region, type |
| VPC | VPC, Subnet, Security Group, Internet/NAT Gateway counts |
| IAM | User and Role counts |
| ECS | Cluster count by region |
| EKS | Cluster count by region |
| SQS | Queue count by region |
| SNS | Topic count by region |
| ACM | Certificate count by region, status |
| Route53 | Hosted zone count |
| ELB | Load balancer count by region, type |

## Quick Start

### Prerequisites

- Docker and Docker Compose
- AWS credentials with read access to the services you want to monitor

### Setup

1. Clone the repository:
```bash
git clone https://github.com/nimishgj/aws-radar.git
cd aws-radar
```

2. Create environment file with AWS credentials:
```bash
cp docker/.env.example docker/.env
# Edit docker/.env with your AWS credentials
```

3. Start all services:
```bash
make docker-up
```

4. Access Grafana dashboard:
- URL: http://localhost:3000
- Username: `admin`
- Password: `admin`

## Configuration

### config.yaml

```yaml
collection:
  interval: 60  # Collection interval in seconds

regions:
  - us-east-1
  - us-west-2
  - eu-west-1
  # Add more regions as needed

collectors:
  - ec2
  - lambda
  - s3
  # Add more collectors as needed
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `AWS_ACCESS_KEY_ID` | AWS access key |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key |
| `AWS_REGION` | Default AWS region (optional) |

## Development

### Build

```bash
make build
```

### Run locally

```bash
make run
```

### Run tests

```bash
make test
```

### Run CI checks

```bash
make ci
```

### Available Make targets

```bash
make help
```

## Project Structure

```
aws-radar/
├── cmd/aws-radar/       # Application entrypoint
├── internal/
│   ├── collector/       # AWS service collectors
│   ├── config/          # Configuration handling
│   ├── metrics/         # Prometheus metrics definitions
│   └── server/          # HTTP server
├── docker/
│   ├── docker-compose.yaml
│   ├── clickhouse/      # ClickHouse init scripts
│   ├── grafana/         # Grafana provisioning
│   └── otel-collector-config.yaml
├── config.yaml          # Default configuration
├── Dockerfile
└── Makefile
```

## Metrics

All metrics are prefixed with `aws_` and have the suffix `_total`. Example metrics:

```
aws_ec2_instances_total{region="us-east-1",instance_type="t3.micro",state="running"} 5
aws_lambda_functions_total{region="us-east-1",runtime="python3.9"} 10
aws_s3_buckets_total{region="us-east-1"} 25
aws_vpc_total{region="us-east-1"} 3
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This project was built by pairing with [Claude Code](https://claude.ai/code), Anthropic's AI-powered coding assistant.
