#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ECS Namespaces"
ALL_NAMESPACE_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Service Discovery Namespaces in $REGION"
  NAMESPACES=$(aws servicediscovery list-namespaces --region "$REGION" \
    --query "Namespaces[].Name" \
    --output text 2>/dev/null || echo "")

  for NAME in $NAMESPACES; do
    ALL_NAMESPACE_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_NAMESPACE_INFO" >> "$OUTPUT_FILE"
