#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="MSK Connectors"
ALL_CONNECTOR_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching MSK Connectors in $REGION"
  CONNECTORS=$(aws kafkaconnect list-connectors --region "$REGION" \
    --query "connectors[].connectorName" \
    --output text 2>/dev/null || echo "")

  for NAME in $CONNECTORS; do
    ALL_CONNECTOR_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_CONNECTOR_INFO" >> "$OUTPUT_FILE"
