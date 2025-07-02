#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch Logs Insights Queries"
ALL_INSIGHTS_QUERIES_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch Logs Insights Saved Queries in $REGION"
  QUERIES=$(aws logs describe-queries --region "$REGION" \
    --status Complete --query "queries[].queryId" \
    --max-items 20 \
    --output text 2>/dev/null || echo "")

  for ID in $QUERIES; do
    ALL_INSIGHTS_QUERIES_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_INSIGHTS_QUERIES_INFO" >> "$OUTPUT_FILE"
