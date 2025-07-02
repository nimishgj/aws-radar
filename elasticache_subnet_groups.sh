#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ElastiCache Subnet Groups"
ALL_SUBNET_GROUP_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ElastiCache Subnet Groups in $REGION"
  SUBNET_GROUPS=$(aws elasticache describe-cache-subnet-groups --region "$REGION" \
    --query "CacheSubnetGroups[].CacheSubnetGroupName" \
    --output text 2>/dev/null || echo "")

  for NAME in $SUBNET_GROUPS; do
    ALL_SUBNET_GROUP_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_SUBNET_GROUP_INFO" >> "$OUTPUT_FILE"
