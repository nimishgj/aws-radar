#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Key Pairs"
ALL_KEYPAIR_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Key Pairs in $REGION"
  KEYPAIR_NAMES=$(aws ec2 describe-key-pairs --region "$REGION" \
    --query "KeyPairs[].KeyName" \
    --output text)

  for NAME in $KEYPAIR_NAMES; do
    ALL_KEYPAIR_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_KEYPAIR_INFO" >> "$OUTPUT_FILE"
