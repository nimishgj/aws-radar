#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ElastiCache Replication Groups"
ALL_GROUP_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ElastiCache Replication Groups in $REGION"
  GROUPS=$(aws elasticache describe-replication-groups --region "$REGION" \
    --query "ReplicationGroups[].ReplicationGroupId" \
    --output text 2>/dev/null || echo "")

  for ID in $GROUPS; do
    ALL_GROUP_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_GROUP_INFO" >> "$OUTPUT_FILE"
