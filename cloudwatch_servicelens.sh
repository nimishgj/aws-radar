#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch ServiceLens/X-Ray Groups"
ALL_SERVICELENS_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching X-Ray Groups for CloudWatch ServiceLens in $REGION"
  GROUPS=$(aws xray get-groups --region "$REGION" \
    --query "Groups[].GroupName" \
    --output text 2>/dev/null || echo "")

  for NAME in $GROUPS; do
    ALL_SERVICELENS_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_SERVICELENS_INFO" >> "$OUTPUT_FILE"
