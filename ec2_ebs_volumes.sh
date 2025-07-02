#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="EBS Volumes"
ALL_VOLUME_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching EBS volumes in $REGION"
  VOLUME_IDS=$(aws ec2 describe-volumes --region "$REGION" \
    --query "Volumes[].VolumeId" \
    --output text)

  for ID in $VOLUME_IDS; do
    ALL_VOLUME_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

# Append a blank line first to separate from previous service entries
echo "" >> "$OUTPUT_FILE"
# Append EBS volumes to the output file
echo "$ALL_VOLUME_INFO" >> "$OUTPUT_FILE"
