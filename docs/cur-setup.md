# CUR (Cost and Usage Report) Setup Guide

AWS Radar can read CUR report files directly from S3 to provide rich cost metrics — per-resource costs, cost by tags, and usage type breakdowns — without incurring the $0.01/call Cost Explorer API charges.

## Prerequisites

- An AWS CUR report configured and delivering to an S3 bucket
- IAM permissions to read the CUR bucket (see below)

## Step 1: Create a CUR Report

If you don't already have a CUR report:

1. Go to the [AWS Billing Console > Cost and Usage Reports](https://console.aws.amazon.com/billing/home#/reports)
2. Click **Create report**
3. Configure the report:
   - **Report name**: e.g., `my-cur-report`
   - **Include resource IDs**: **Yes** (required for per-resource cost metrics)
   - **Data refresh settings**: Enable automatic refresh
4. Choose delivery options:
   - **S3 bucket**: Select or create a bucket (e.g., `my-cur-bucket`)
   - **Report path prefix**: e.g., `cur-reports`
   - **Time granularity**: Hourly or Daily
   - **Report versioning**: Overwrite existing report
   - **Compression type**: **Parquet** (recommended) or GZIP
5. Click **Review and Complete**

For detailed instructions, see the [AWS documentation](https://docs.aws.amazon.com/cur/latest/userguide/cur-create.html).

### Recommended Settings

- **Format**: Parquet (smaller files, faster to read)
- **Include resource IDs**: Must be enabled for `aws_cur_cost_by_resource_usd`
- **Include user-defined cost allocation tags**: Enable in Billing > Cost Allocation Tags for `aws_cur_cost_by_tag_usd`

## Step 2: Configure aws-radar

Add the `cost_cur` section to your `config.yaml`:

```yaml
cost_cur:
  enabled: true
  bucket: my-cur-bucket
  prefix: cur-reports
  report_name: my-cur-report
  format: ""           # auto-detect (parquet or csv); or set explicitly
  frequency: daily     # how often to re-read CUR (hourly | daily)
  region: us-east-1    # S3 bucket region
  max_resources: 100   # top N resources by cost to emit as metrics
```

## Step 3: IAM Permissions

Add the following permissions to your aws-radar IAM policy:

```json
{
    "Sid": "AllowCURBucketRead",
    "Effect": "Allow",
    "Action": [
        "s3:GetObject",
        "s3:ListBucket"
    ],
    "Resource": [
        "arn:aws:s3:::my-cur-bucket",
        "arn:aws:s3:::my-cur-bucket/*"
    ]
}
```

Replace `my-cur-bucket` with your actual CUR bucket name.

## Metrics Produced

| Metric | Labels | Description |
|--------|--------|-------------|
| `aws_cur_total_cost_usd` | account, account_name, period | Total cost for the billing period |
| `aws_cur_cost_by_service_usd` | account, account_name, service, period | Cost broken down by AWS service |
| `aws_cur_cost_by_resource_usd` | account, account_name, service, resource_id, period | Top N resources by cost |
| `aws_cur_cost_by_usage_type_usd` | account, account_name, service, usage_type, period | Cost by usage type |
| `aws_cur_cost_by_tag_usd` | account, account_name, tag_key, tag_value, period | Cost by user-defined tag |
| `aws_cur_last_processed_timestamp` | account, account_name | Unix timestamp of last successful processing |

## How It Works

1. AWS Radar computes the current billing period path: `<prefix>/<report-name>/yyyymmdd-yyyymmdd/`
2. Downloads `<report-name>-Manifest.json` from that path
3. Reads the `reportKeys` from the manifest (list of data file S3 keys)
4. For each data file:
   - **Parquet**: Downloads to a temp file, reads row-by-row
   - **CSV (gzip)**: Streams via S3 GetObject with gzip decompression
5. Aggregates costs into maps: by service, by resource (top N), by usage type, by tag
6. Emits Prometheus metrics
7. Caches results; refreshes based on configured frequency

## CUR vs Cost Explorer

| Feature | Cost Explorer | CUR |
|---------|--------------|-----|
| API cost | $0.01/call | Free (S3 storage only) |
| Granularity | Service-level | Per-resource, per-tag |
| Data freshness | Near real-time | Updated ~1x/day |
| Setup | Enable in config | Create CUR report + S3 bucket |

Both can be enabled simultaneously. CUR provides richer data; Cost Explorer provides fresher data.
