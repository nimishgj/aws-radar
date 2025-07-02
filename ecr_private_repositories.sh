#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ECR Private Repositories"
ALL_REPO_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ECR Private Repositories in $REGION"
  REPO_NAMES=$(aws ecr describe-repositories --region "$REGION" \
    --query "repositories[].repositoryName" \
    --output text 2>/dev/null || echo "")

  for NAME in $REPO_NAMES; do
    ALL_REPO_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_REPO_INFO" >> "$OUTPUT_FILE"
