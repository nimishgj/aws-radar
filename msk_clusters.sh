#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="MSK Clusters"
ALL_CLUSTER_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching MSK Clusters in $REGION"
  CLUSTERS=$(aws kafka list-clusters --region "$REGION" \
    --query "ClusterInfoList[].ClusterName" \
    --output text 2>/dev/null || echo "")

  for NAME in $CLUSTERS; do
    ALL_CLUSTER_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_CLUSTER_INFO" >> "$OUTPUT_FILE"
