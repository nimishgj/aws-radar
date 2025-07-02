#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch Evidently Projects"
ALL_EVIDENTLY_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch Evidently Projects in $REGION"
  PROJECTS=$(aws evidently list-projects --region "$REGION" \
    --query "projects[].name" \
    --output text 2>/dev/null || echo "")

  for NAME in $PROJECTS; do
    ALL_EVIDENTLY_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_EVIDENTLY_INFO" >> "$OUTPUT_FILE"
