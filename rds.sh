#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="RDS Instances"
ALL_RDS_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching RDS Instances in $REGION"
  RDS_INSTANCES=$(aws rds describe-db-instances --region "$REGION" \
    --query "DBInstances[].DBInstanceIdentifier" \
    --output text 2>/dev/null || echo "")

  for ID in $RDS_INSTANCES; do
    ALL_RDS_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_RDS_INFO" >> "$OUTPUT_FILE"
