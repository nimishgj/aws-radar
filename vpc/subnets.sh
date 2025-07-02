#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Subnets"
ALL_SUBNET_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Subnets in $REGION"
  SUBNET_IDS=$(aws ec2 describe-subnets --region "$REGION" \
    --query "Subnets[].SubnetId" \
    --output text)

  for ID in $SUBNET_IDS; do
    ALL_SUBNET_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_SUBNET_INFO" >> "$OUTPUT_FILE"
