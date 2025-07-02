#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Security Groups"
ALL_SG_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching security groups in $REGION"
  SG_IDS=$(aws ec2 describe-security-groups --region "$REGION" \
    --query "SecurityGroups[].GroupId" \
    --output text)

  for ID in $SG_IDS; do
    ALL_SG_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_SG_INFO" >> "$OUTPUT_FILE"
