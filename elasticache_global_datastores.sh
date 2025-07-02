#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ElastiCache Global Datastores"
ALL_GLOBAL_DATASTORE_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ElastiCache Global Datastores in $REGION"
  # Get global datastores (only available with Redis and Valkey)
  GLOBAL_DATASTORES=$(aws elasticache describe-global-replication-groups --region "$REGION" \
    --query "GlobalReplicationGroups[].GlobalReplicationGroupId" \
    --output text 2>/dev/null || echo "")

  for ID in $GLOBAL_DATASTORES; do
    ALL_GLOBAL_DATASTORE_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_GLOBAL_DATASTORE_INFO" >> "$OUTPUT_FILE"
