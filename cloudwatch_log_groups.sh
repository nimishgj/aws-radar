#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch Log Groups"
ALL_LOG_GROUP_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch Log Groups in $REGION"
  LOG_GROUPS=$(aws logs describe-log-groups --region "$REGION" \
    --query "logGroups[].logGroupName" \
    --output text 2>/dev/null || echo "")

  for NAME in $LOG_GROUPS; do
    ALL_LOG_GROUP_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_LOG_GROUP_INFO" >> "$OUTPUT_FILE"
