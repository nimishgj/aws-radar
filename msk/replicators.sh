#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="MSK Replicators"
ALL_REPLICATOR_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching MSK Replicators in $REGION"
  REPLICATORS=$(aws kafka list-replicators --region "$REGION" \
    --query "replicators[].replicatorName" \
    --output text 2>/dev/null || echo "")

  for NAME in $REPLICATORS; do
    ALL_REPLICATOR_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_REPLICATOR_INFO" >> "$OUTPUT_FILE"
