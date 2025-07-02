#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="KMS Keys"
ALL_KEY_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching KMS Keys in $REGION"
  KEYS=$(aws kms list-keys --region "$REGION" \
    --query "Keys[].KeyId" \
    --output text 2>/dev/null || echo "")

  for KEY_ID in $KEYS; do
    # Try to get the key alias if available for better identification
    KEY_ALIAS=$(aws kms list-aliases --region "$REGION" \
      --key-id "$KEY_ID" \
      --query "Aliases[].AliasName" \
      --output text 2>/dev/null || echo "")
    
    if [ -n "$KEY_ALIAS" ]; then
      KEY_IDENTIFIER="$KEY_ALIAS ($KEY_ID)"
    else
      KEY_IDENTIFIER="$KEY_ID"
    fi
    
    ALL_KEY_INFO+="$RESOURCE_TYPE,$KEY_IDENTIFIER,$REGION"$'\n'
  done
done

echo "$ALL_KEY_INFO" >> "$OUTPUT_FILE"
