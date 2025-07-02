#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Secrets Manager Secrets"
ALL_SECRET_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Secrets Manager Secrets in $REGION"
  SECRETS=$(aws secretsmanager list-secrets --region "$REGION" \
    --query "SecretList[].Name" \
    --output text 2>/dev/null || echo "")

  for SECRET_NAME in $SECRETS; do
    ALL_SECRET_INFO+="$RESOURCE_TYPE,$SECRET_NAME,$REGION"$'\n'
  done
done

echo "$ALL_SECRET_INFO" >> "$OUTPUT_FILE"
