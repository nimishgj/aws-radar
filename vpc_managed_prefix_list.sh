#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Managed Prefix Lists"
ALL_PREFIX_LIST_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Managed Prefix Lists in $REGION"
  PREFIX_LIST_IDS=$(aws ec2 describe-managed-prefix-lists --region "$REGION" \
    --query "PrefixLists[].PrefixListId" \
    --output text 2>/dev/null || echo "")

  for ID in $PREFIX_LIST_IDS; do
    ALL_PREFIX_LIST_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_PREFIX_LIST_INFO" >> "$OUTPUT_FILE"
