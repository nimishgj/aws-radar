#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="EKS Clusters"
ALL_CLUSTER_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching EKS Clusters in $REGION"
  CLUSTERS=$(aws eks list-clusters --region "$REGION" \
    --query "clusters[]" \
    --output text 2>/dev/null || echo "")

  for CLUSTER in $CLUSTERS; do
    ALL_CLUSTER_INFO+="$RESOURCE_TYPE,$CLUSTER,$REGION"$'\n'
  done
done

echo "$ALL_CLUSTER_INFO" >> "$OUTPUT_FILE"
