#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Route Tables"
ALL_RT_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Route Tables in $REGION"
  RT_IDS=$(aws ec2 describe-route-tables --region "$REGION" \
    --query "RouteTables[].RouteTableId" \
    --output text)

  for ID in $RT_IDS; do
    ALL_RT_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_RT_INFO" >> "$OUTPUT_FILE"
