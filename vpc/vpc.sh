#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC"
ALL_VPC_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPCs in $REGION"
  VPC_IDS=$(aws ec2 describe-vpcs --region "$REGION" \
    --query "Vpcs[].VpcId" \
    --output text)

  for ID in $VPC_IDS; do
    ALL_VPC_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_VPC_INFO" >> "$OUTPUT_FILE"
