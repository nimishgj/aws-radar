#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="ECS Task Definitions"
ALL_TASKDEF_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching ECS Task Definitions in $REGION"
  TASKDEF_ARNS=$(aws ecs list-task-definitions --region "$REGION" \
    --query "taskDefinitionArns[]" \
    --output text 2>/dev/null || echo "")

  for ARN in $TASKDEF_ARNS; do
    # Extract task definition name from ARN
    TASKDEF_NAME=$(echo $ARN | awk -F'/' '{print $2}')
    ALL_TASKDEF_INFO+="$RESOURCE_TYPE,$TASKDEF_NAME,$REGION"$'\n'
  done
done

echo "$ALL_TASKDEF_INFO" >> "$OUTPUT_FILE"
