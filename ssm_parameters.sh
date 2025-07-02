#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="SSM Parameters"
ALL_PARAM_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching SSM Parameters in $REGION"
  # Use describe-parameters because get-parameters requires knowing the names
  PARAMS=$(aws ssm describe-parameters --region "$REGION" \
    --query "Parameters[].Name" \
    --output text 2>/dev/null || echo "")

  for PARAM_NAME in $PARAMS; do
    ALL_PARAM_INFO+="$RESOURCE_TYPE,$PARAM_NAME,$REGION"$'\n'
  done
done

echo "$ALL_PARAM_INFO" >> "$OUTPUT_FILE"
