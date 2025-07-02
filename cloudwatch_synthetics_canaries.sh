#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch Synthetics Canaries"
ALL_CANARY_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch Synthetics Canaries in $REGION"
  CANARIES=$(aws synthetics describe-canaries --region "$REGION" \
    --query "Canaries[].Name" \
    --output text 2>/dev/null || echo "")

  for NAME in $CANARIES; do
    ALL_CANARY_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_CANARY_INFO" >> "$OUTPUT_FILE"
