#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="RDS Snapshots"
ALL_SNAPSHOT_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching RDS Snapshots in $REGION"
  SNAPSHOTS=$(aws rds describe-db-snapshots --region "$REGION" \
    --query "DBSnapshots[].DBSnapshotIdentifier" \
    --output text 2>/dev/null || echo "")

  for ID in $SNAPSHOTS; do
    ALL_SNAPSHOT_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_SNAPSHOT_INFO" >> "$OUTPUT_FILE"
