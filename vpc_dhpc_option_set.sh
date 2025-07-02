#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC DHCP Options"
ALL_DHCP_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC DHCP Option Sets in $REGION"
  DHCP_IDS=$(aws ec2 describe-dhcp-options --region "$REGION" \
    --query "DhcpOptions[].DhcpOptionsId" \
    --output text)

  for ID in $DHCP_IDS; do
    ALL_DHCP_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_DHCP_INFO" >> "$OUTPUT_FILE"
