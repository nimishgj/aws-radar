#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ElastiCache Memcached Caches"
ALL_MEMCACHED_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ElastiCache Memcached Caches in $REGION"
  # Filter clusters with memcached engine
  MEMCACHED_CLUSTERS=$(aws elasticache describe-cache-clusters --region "$REGION" \
    --query "CacheClusters[?Engine=='memcached'].CacheClusterId" \
    --output text 2>/dev/null || echo "")

  for ID in $MEMCACHED_CLUSTERS; do
    ALL_MEMCACHED_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_MEMCACHED_INFO" >> "$OUTPUT_FILE"
