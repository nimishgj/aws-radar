#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="API Gateway HTTP APIs"
ALL_API_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching API Gateway HTTP APIs in $REGION"
  API_IDS=$(aws apigatewayv2 get-apis --region "$REGION" \
    --query "Items[?ApiType=='HTTP'].[ApiId,Name]" \
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
