#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Transit Gateways"
ALL_TGW_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Transit Gateways in $REGION"
  TGW_IDS=$(aws ec2 describe-transit-gateways --region "$REGION" \
    --query "TransitGateways[].TransitGatewayId" \
    --output text 2>/dev/null || echo "")

  for ID in $TGW_IDS; do
    ALL_TGW_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_TGW_INFO" >> "$OUTPUT_FILE"
