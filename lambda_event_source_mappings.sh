#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Lambda Event Source Mappings"
ALL_MAPPING_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Lambda Event Source Mappings in $REGION"
  MAPPINGS=$(aws lambda list-event-source-mappings --region "$REGION" \
    --query "EventSourceMappings[].UUID" \
    --output text 2>/dev/null || echo "")

  for UUID in $MAPPINGS; do
    ALL_MAPPING_INFO+="$RESOURCE_TYPE,$UUID,$REGION"$'\n'
  done
done

echo "$ALL_MAPPING_INFO" >> "$OUTPUT_FILE"
