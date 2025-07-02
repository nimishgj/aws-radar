#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ElastiCache Security Groups"
ALL_SG_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ElastiCache Security Groups in $REGION"
  SECURITY_GROUPS=$(aws elasticache describe-cache-security-groups --region "$REGION" \
    --query "CacheSecurityGroups[].CacheSecurityGroupName" \
    --output text 2>/dev/null || echo "")

  for NAME in $SECURITY_GROUPS; do
    ALL_SG_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_SG_INFO" >> "$OUTPUT_FILE"
