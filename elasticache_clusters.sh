#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ElastiCache Clusters"
ALL_CLUSTER_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ElastiCache Clusters in $REGION"
  CLUSTERS=$(aws elasticache describe-cache-clusters --region "$REGION" \
    --query "CacheClusters[].CacheClusterId" \
    --output text 2>/dev/null || echo "")

  for ID in $CLUSTERS; do
    ALL_CLUSTER_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_CLUSTER_INFO" >> "$OUTPUT_FILE"
