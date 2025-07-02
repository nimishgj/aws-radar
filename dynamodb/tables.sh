#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="DynamoDB Tables"
ALL_TABLE_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching DynamoDB Tables in $REGION"
  TABLES=$(aws dynamodb list-tables --region "$REGION" \
    --query "TableNames[]" \
    --output text 2>/dev/null || echo "")

  for TABLE_NAME in $TABLES; do
    ALL_TABLE_INFO+="$RESOURCE_TYPE,$TABLE_NAME,$REGION"$'\n'
  done
done

echo "$ALL_TABLE_INFO" >> "$OUTPUT_FILE"
