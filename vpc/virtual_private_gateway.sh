#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Virtual Private Gateways"
ALL_VPG_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Virtual Private Gateways in $REGION"
  VPG_IDS=$(aws ec2 describe-vpn-gateways --region "$REGION" \
    --query "VpnGateways[].VpnGatewayId" \
    --output text)

  for ID in $VPG_IDS; do
    ALL_VPG_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_VPG_INFO" >> "$OUTPUT_FILE"
