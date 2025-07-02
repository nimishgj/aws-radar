#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Auto Scaling Groups"
ALL_ASG_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Auto Scaling Groups in $REGION"
  ASG_NAMES=$(aws autoscaling describe-auto-scaling-groups --region "$REGION" \
    --query "AutoScalingGroups[].AutoScalingGroupName" \
    --output text 2>/dev/null || echo "")

  for NAME in $ASG_NAMES; do
    ALL_ASG_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_ASG_INFO" >> "$OUTPUT_FILE"
