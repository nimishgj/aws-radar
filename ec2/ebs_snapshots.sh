#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="EBS Snapshots"
ALL_SNAPSHOT_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching EBS snapshots in $REGION"
  SNAPSHOT_IDS=$(aws ec2 describe-snapshots --owner-ids self --region "$REGION" \
    --query "Snapshots[].SnapshotId" \
    --output text)

  for ID in $SNAPSHOT_IDS; do
    ALL_SNAPSHOT_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_SNAPSHOT_INFO" >> "$OUTPUT_FILE"
