#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ElastiCache Redis OSS Caches"
ALL_REDIS_OSS_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ElastiCache Redis OSS Caches in $REGION"
  # Filter clusters with redis engine
  REDIS_CLUSTERS=$(aws elasticache describe-cache-clusters --region "$REGION" \
    --query "CacheClusters[?Engine=='redis'].CacheClusterId" \
    --output text 2>/dev/null || echo "")

  for ID in $REDIS_CLUSTERS; do
    ALL_REDIS_OSS_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_REDIS_OSS_INFO" >> "$OUTPUT_FILE"
