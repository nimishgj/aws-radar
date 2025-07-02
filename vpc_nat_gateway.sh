#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC NAT Gateways"
ALL_NAT_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC NAT Gateways in $REGION"
  NAT_IDS=$(aws ec2 describe-nat-gateways --region "$REGION" \
    --query "NatGateways[].NatGatewayId" \
    --output text)

  for ID in $NAT_IDS; do
    ALL_NAT_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_NAT_INFO" >> "$OUTPUT_FILE"
