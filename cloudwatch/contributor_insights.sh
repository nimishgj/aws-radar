#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch Contributor Insights"
ALL_INSIGHTS_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch Contributor Insights in $REGION"
  RULES=$(aws cloudwatch list-insight-rules --region "$REGION" \
    --query "InsightRules[].Name" \
    --output text 2>/dev/null || echo "")

  for NAME in $RULES; do
    ALL_INSIGHTS_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_INSIGHTS_INFO" >> "$OUTPUT_FILE"
