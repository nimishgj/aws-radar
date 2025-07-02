#!/bin/bash

# Output file name
OUTPUT_FILE="aws_resources.csv"

# Initialize the file
> "$OUTPUT_FILE"

# Global: Get all AWS regions
echo "INFO: fetching aws regions"
AWS_REGIONS=($(aws ec2 describe-regions --query "Regions[].RegionName" --output text))

# --- S3 ---
RESOURCE_TYPE="S3"
S3_BUCKET_NAMES=$(aws s3 ls | awk '{print $3}')

for BUCKET in $S3_BUCKET_NAMES; do
  echo "$RESOURCE_TYPE,$BUCKET,ALL_REGION" >> "$OUTPUT_FILE"
done

echo "INFO: fetching s3 buckets"

# --- EC2 ---
RESOURCE_TYPE="EC2"
ALL_INSTANCE_NAMES=""

for REGION in "${AWS_REGIONS[@]}"; do
  echo "INFO: fetching ec2 instances in $REGION"
  INSTANCE_NAMES=$(aws ec2 describe-instances --region "$REGION" \
    --query "Reservations[].Instances[].Tags[?Key=='Name'].Value | []" \
    --output text)

  for NAME in $INSTANCE_NAMES; do
    ALL_INSTANCE_NAMES+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

# Append EC2 names to the output file
echo "$ALL_INSTANCE_NAMES" >> "$OUTPUT_FILE"
echo "INFO: output $OUTPUT_FILE"

