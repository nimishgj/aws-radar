#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Step Functions State Machines"
ALL_SM_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Step Functions State Machines in $REGION"
  STATE_MACHINES=$(aws stepfunctions list-state-machines --region "$REGION" \
    --query "stateMachines[].[name,stateMachineArn]" \
    --output text 2>/dev/null || echo "")

  if [ -n "$STATE_MACHINES" ]; then
    while read -r SM_NAME SM_ARN; do
      ALL_SM_INFO+="$RESOURCE_TYPE,$SM_NAME,$REGION"$'\n'
    done <<< "$STATE_MACHINES"
  fi
done

echo "$ALL_SM_INFO" >> "$OUTPUT_FILE"
