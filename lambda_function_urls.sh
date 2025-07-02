#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Lambda Function URLs"
ALL_URL_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Lambda Function URLs in $REGION"
  # First get all functions
  FUNCTIONS=$(aws lambda list-functions --region "$REGION" \
    --query "Functions[].FunctionName" \
    --output text 2>/dev/null || echo "")

  for FUNC in $FUNCTIONS; do
    # Check if function has a URL configuration
    URL_CONFIG=$(aws lambda list-function-url-configs --function-name "$FUNC" \
      --region "$REGION" --query "FunctionUrlConfigs[].FunctionUrl" \
      --output text 2>/dev/null || echo "")
    
    if [ ! -z "$URL_CONFIG" ]; then
      ALL_URL_INFO+="$RESOURCE_TYPE,$FUNC,$REGION"$'\n'
    fi
  done
done

echo "$ALL_URL_INFO" >> "$OUTPUT_FILE"
