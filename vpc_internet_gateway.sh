#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Internet Gateways"
ALL_IGW_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Internet Gateways in $REGION"
  IGW_IDS=$(aws ec2 describe-internet-gateways --region "$REGION" \
    --query "InternetGateways[].InternetGatewayId" \
    --output text)

  for ID in $IGW_IDS; do
    ALL_IGW_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_IGW_INFO" >> "$OUTPUT_FILE"
