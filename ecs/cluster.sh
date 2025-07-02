#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ECS Clusters"
ALL_CLUSTER_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ECS Clusters in $REGION"
  CLUSTER_ARNS=$(aws ecs list-clusters --region "$REGION" \
    --query "clusterArns[]" \
    --output text 2>/dev/null || echo "")

  for ARN in $CLUSTER_ARNS; do
    # Extract cluster name from ARN
    CLUSTER_NAME=$(echo $ARN | awk -F'/' '{print $2}')
    ALL_CLUSTER_INFO+="$RESOURCE_TYPE,$CLUSTER_NAME,$REGION"$'\n'
  done
done

echo "$ALL_CLUSTER_INFO" >> "$OUTPUT_FILE"
