# AWS Radar - Project Requirements & Status

## Overview

AWS Radar is a Go-based monitoring agent that collects AWS resource metrics and exports them in Prometheus format. An OpenTelemetry Collector scrapes these metrics and forwards them to a backend.

## Architecture

```
┌─────────────┐      ┌─────────────────────┐      ┌─────────────────┐      ┌─────────────┐
│   AWS APIs  │ ──>  │   aws-radar         │ <──  │ OTel Collector  │ ──>  │  Backend    │
│             │      │  (:9090/metrics)    │ pull │                 │      │ (ClickHouse │
└─────────────┘      └─────────────────────┘      └─────────────────┘      │  or OTLP)   │
                            ↑                            │                  └─────────────┘
                     Prometheus format            Scrapes every 60s
```

## Current Implementation Status

### Completed

| Component | Status | Details |
|-----------|--------|---------|
| AWS API Integration | Done | 25 collectors across 18 services |
| Prometheus `/metrics` endpoint | Done | Exposes all metrics |
| Multi-region support | Done | Configurable via config.yaml |
| Parallel collection | Done | WaitGroups for concurrent execution |
| Docker Compose stack | Done | aws-radar + otel-collector + clickhouse + grafana |
| Grafana dashboards | Done | Pre-configured ClickHouse datasource |
| ClickHouse schema | Done | With views and TTL |

### AWS Services Covered

**Regional Collectors:**
- EC2 (instances, by type/state/AZ)
- S3 (buckets)
- RDS (instances, by class/engine/multi-AZ)
- Lambda (functions, by runtime/memory)
- ECS (services, tasks)
- EKS (clusters)
- DynamoDB (tables)
- ElastiCache (clusters)
- SQS (queues)
- SNS (topics)
- EBS (volumes)
- VPC (VPCs, subnets, security groups, NAT/Internet gateways)
- ACM (certificates)
- ELB/ALB/NLB (load balancers)

**Global Collectors:**
- CloudFront (distributions)
- Route53 (hosted zones)
- IAM (users, roles)

## Gaps & Required Work

### Critical (Must Have)

| Item | Priority | Description |
|------|----------|-------------|
| Unit Tests | High | 0 test files currently - need tests for collectors |
| README | High | No documentation for setup/usage |
| Build Verification | High | Verify compilation and runtime |
| IAM Permissions Doc | High | Document required AWS permissions |

### Important (Should Have)

| Item | Priority | Description |
|------|----------|-------------|
| Retry Logic | Medium | AWS API calls can fail transiently |
| Rate Limiting | Medium | Prevent AWS API throttling |
| Error Recovery | Medium | Failed collectors don't retry until next cycle |
| Pagination Verification | Medium | Ensure large result sets handled correctly |
| Metrics Validation | Medium | Verify metrics appear in backend |

### Nice to Have

| Item | Priority | Description |
|------|----------|-------------|
| CI/CD Pipeline | Low | Automated testing and builds |
| Linting Setup | Low | golangci-lint configuration |
| Contributing Guide | Low | For future contributors |
| Helm Chart | Low | Kubernetes deployment |
| Alerting Rules | Low | Prometheus/Grafana alerts |

## Project Structure

```
aws-radar/
├── cmd/aws-radar/main.go           # Entry point
├── internal/
│   ├── config/config.go            # Viper-based config
│   ├── metrics/metrics.go          # Prometheus metrics
│   ├── collector/                  # AWS collectors (18 files)
│   └── server/server.go            # HTTP server
├── docker/                         # Docker-related configs
│   ├── docker-compose.yaml
│   ├── otel-collector-config.yaml
│   ├── infraspec.env
│   ├── clickhouse/init.sql
│   └── grafana/
├── docs/                           # Documentation
├── config.yaml                     # App configuration
├── Dockerfile
├── Makefile
└── go.mod
```

## Configuration

### Environment Variables

| Variable | Description |
|----------|-------------|
| AWS_ACCESS_KEY_ID | AWS access key |
| AWS_SECRET_ACCESS_KEY | AWS secret key |
| AWS_RADAR_CONFIG | Path to config file |

### config.yaml

```yaml
server:
  port: 9090
  metrics_path: /metrics
  health_path: /health

collection:
  interval: 60s
  timeout: 30s

aws:
  regions:
    - us-east-1
    - us-west-2
    - eu-west-1

logging:
  level: info
  format: json
```

## Running the Application

### Local Development

```bash
make build
./bin/aws-radar
```

### Docker Compose

```bash
cd docker
docker-compose up -d
```

### Endpoints

- `http://localhost:9090/metrics` - Prometheus metrics
- `http://localhost:9090/health` - Health check
- `http://localhost:3000` - Grafana (admin/admin)
- `http://localhost:8123` - ClickHouse HTTP

## Dependencies

### Go Modules
- aws-sdk-go-v2 (15+ AWS services)
- prometheus/client_golang
- spf13/viper
- rs/zerolog

### Infrastructure
- OpenTelemetry Collector Contrib
- ClickHouse
- Grafana

## Next Steps

1. Verify build: `make build`
2. Run locally with AWS credentials
3. Add unit tests for core collectors
4. Create README with quickstart guide
5. Document IAM permissions required
