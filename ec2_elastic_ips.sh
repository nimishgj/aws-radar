#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Elastic IPs"
ALL_EIP_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Elastic IPs in $REGION"
  EIP_ADDRESSES=$(aws ec2 describe-addresses --region "$REGION" \
    --query "Addresses[].PublicIp" \
    --output text)

  for IP in $EIP_ADDRESSES; do
    ALL_EIP_INFO+="$RESOURCE_TYPE,$IP,$REGION"$'\n'
  done
done

echo "$ALL_EIP_INFO" >> "$OUTPUT_FILE"
