#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ECR Public Repositories"
ALL_REPO_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ECR Public Repositories in $REGION"
  
  # Note: Public ECR is only available in us-east-1
  if [ "$REGION" == "us-east-1" ]; then
    REPO_NAMES=$(aws ecr-public describe-repositories --region "$REGION" \
      --query "repositories[].repositoryName" \
      --output text 2>/dev/null || echo "")

    for NAME in $REPO_NAMES; do
      ALL_REPO_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
    done
  fi
done

echo "$ALL_REPO_INFO" >> "$OUTPUT_FILE"
