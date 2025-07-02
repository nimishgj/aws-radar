#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="VPC Network ACLs"
ALL_NACL_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching VPC Network ACLs in $REGION"
  NACL_IDS=$(aws ec2 describe-network-acls --region "$REGION" \
    --query "NetworkAcls[].NetworkAclId" \
    --output text)

  for ID in $NACL_IDS; do
    ALL_NACL_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
  done
done

echo "$ALL_NACL_INFO" >> "$OUTPUT_FILE"
