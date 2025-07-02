#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Customer Gateways"
ALL_CGW_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Customer Gateways in $REGION"
  CGW_IDS=$(aws ec2 describe-customer-gateways --region "$REGION" \
    --query "CustomerGateways[].CustomerGatewayId" \
    --output text)

  for ID in $CGW_IDS; do
    ALL_CGW_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_CGW_INFO" >> "$OUTPUT_FILE"
