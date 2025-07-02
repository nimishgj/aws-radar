#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Network Interfaces"
ALL_ENI_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Network Interfaces in $REGION"
  ENI_IDS=$(aws ec2 describe-network-interfaces --region "$REGION" \
    --query "NetworkInterfaces[].NetworkInterfaceId" \
    --output text)

  for ID in $ENI_IDS; do
    ALL_ENI_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_ENI_INFO" >> "$OUTPUT_FILE"
