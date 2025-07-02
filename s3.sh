#!/bin/bash

OUTPUT_FILE="$1"
RESOURCE_TYPE="S3"

echo "INFO: Fetching S3 buckets"

S3_BUCKET_NAMES=$(aws s3 ls | awk '{print $3}')

for BUCKET in $S3_BUCKET_NAMES; do
  echo "$RESOURCE_TYPE,$BUCKET,ALL_REGION" >> "$OUTPUT_FILE"
done

