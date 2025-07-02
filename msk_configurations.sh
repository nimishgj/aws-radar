#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="MSK Configurations"
ALL_CONFIG_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching MSK Configurations in $REGION"
  CONFIGS=$(aws kafka list-configurations --region "$REGION" \
    --query "Configurations[].Name" \
    --output text 2>/dev/null || echo "")

  for NAME in $CONFIGS; do
    ALL_CONFIG_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_CONFIG_INFO" >> "$OUTPUT_FILE"
