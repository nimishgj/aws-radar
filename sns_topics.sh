#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="SNS Topics"
ALL_TOPIC_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching SNS Topics in $REGION"
  TOPICS=$(aws sns list-topics --region "$REGION" \
    --query "Topics[].TopicArn" \
    --output text 2>/dev/null || echo "")

  for TOPIC_ARN in $TOPICS; do
    # Extract the topic name from the ARN
    TOPIC_NAME=$(echo "$TOPIC_ARN" | awk -F':' '{print $NF}')
    ALL_TOPIC_INFO+="$RESOURCE_TYPE,$TOPIC_NAME,$REGION"$'\n'
  done
done

echo "$ALL_TOPIC_INFO" >> "$OUTPUT_FILE"
