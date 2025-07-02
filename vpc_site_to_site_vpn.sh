#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Site-to-Site VPN"
ALL_VPN_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Site-to-Site VPN Connections in $REGION"
  VPN_IDS=$(aws ec2 describe-vpn-connections --region "$REGION" \
    --query "VpnConnections[].VpnConnectionId" \
    --output text)

  for ID in $VPN_IDS; do
    ALL_VPN_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_VPN_INFO" >> "$OUTPUT_FILE"
