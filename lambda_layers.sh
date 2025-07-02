#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Lambda Layers"
ALL_LAYER_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Lambda Layers in $REGION"
  LAYERS=$(aws lambda list-layers --region "$REGION" \
    --query "Layers[].LayerName" \
    --output text 2>/dev/null || echo "")

  for NAME in $LAYERS; do
    ALL_LAYER_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_LAYER_INFO" >> "$OUTPUT_FILE"
