#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Target Groups"
ALL_TG_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Target Groups in $REGION"
  TG_ARNS=$(aws elbv2 describe-target-groups --region "$REGION" \
    --query "TargetGroups[].TargetGroupName" \
    --output text 2>/dev/null || echo "")

  for NAME in $TG_ARNS; do
    ALL_TG_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_TG_INFO" >> "$OUTPUT_FILE"
