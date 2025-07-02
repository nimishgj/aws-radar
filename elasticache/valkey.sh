#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ElastiCache Valkey Caches"
ALL_VALKEY_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ElastiCache Valkey Caches in $REGION"
  # Filter clusters with valkey engine
  VALKEY_CLUSTERS=$(aws elasticache describe-cache-clusters --region "$REGION" \
    --query "CacheClusters[?Engine=='valkey'].CacheClusterId" \
    --output text 2>/dev/null || echo "")

  for ID in $VALKEY_CLUSTERS; do
    ALL_VALKEY_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_VALKEY_INFO" >> "$OUTPUT_FILE"
