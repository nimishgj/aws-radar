#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Egress-Only IGW"
ALL_EIGW_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Egress-Only Internet Gateways in $REGION"
  EIGW_IDS=$(aws ec2 describe-egress-only-internet-gateways --region "$REGION" \
    --query "EgressOnlyInternetGateways[].EgressOnlyInternetGatewayId" \
    --output text)

  for ID in $EIGW_IDS; do
    ALL_EIGW_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_EIGW_INFO" >> "$OUTPUT_FILE"
