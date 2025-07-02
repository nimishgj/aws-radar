#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="API Gateway REST APIs"
ALL_API_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching API Gateway REST APIs in $REGION"
  API_IDS=$(aws apigateway get-rest-apis --region "$REGION" \
    --query "items[*].[id,name]" \
    --output text 2>/dev/null || echo "")

  if [ -n "$API_IDS" ]; then
    while read -r ID NAME; do
      if [ -n "$ID" ] && [ -n "$NAME" ]; then
        ALL_API_INFO+="$RESOURCE_TYPE,$NAME ($ID),$REGION"$'\n'
      fi
    done <<< "$API_IDS"
  fi
done

echo "$ALL_API_INFO" >> "$OUTPUT_FILE"
