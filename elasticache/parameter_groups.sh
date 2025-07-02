#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ElastiCache Parameter Groups"
ALL_PARAM_GROUP_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ElastiCache Parameter Groups in $REGION"
  PARAM_GROUPS=$(aws elasticache describe-cache-parameter-groups --region "$REGION" \
    --query "CacheParameterGroups[].CacheParameterGroupName" \
    --output text 2>/dev/null || echo "")

  for NAME in $PARAM_GROUPS; do
    ALL_PARAM_GROUP_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_PARAM_GROUP_INFO" >> "$OUTPUT_FILE"
