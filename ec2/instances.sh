#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="EC2"
ALL_INSTANCE_NAMES=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching EC2 instances in $REGION"
  INSTANCE_NAMES=$(aws ec2 describe-instances --region "$REGION" \
    --query "Reservations[].Instances[].Tags[?Key=='Name'].Value | []" \
    --output text)

  for NAME in $INSTANCE_NAMES; do
    ALL_INSTANCE_NAMES+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

# Append EC2 instances to the output file
echo "$ALL_INSTANCE_NAMES" >> "$OUTPUT_FILE"

