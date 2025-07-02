#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Lambda Functions"
ALL_FUNCTION_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Lambda Functions in $REGION"
  FUNCTIONS=$(aws lambda list-functions --region "$REGION" \
    --query "Functions[].FunctionName" \
    --output text 2>/dev/null || echo "")

  for NAME in $FUNCTIONS; do
    ALL_FUNCTION_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_FUNCTION_INFO" >> "$OUTPUT_FILE"
