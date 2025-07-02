#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="SQS Queues"
ALL_QUEUE_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching SQS Queues in $REGION"
  QUEUE_URLS=$(aws sqs list-queues --region "$REGION" \
    --query "QueueUrls[]" \
    --output text 2>/dev/null || echo "")
  
  if [ -n "$QUEUE_URLS" ]; then
    while IFS= read -r QUEUE_URL; do
      # Extract queue name from URL
      QUEUE_NAME=$(echo "$QUEUE_URL" | awk -F'/' '{print $NF}')
      
      # Check if it's a FIFO queue
      if [[ "$QUEUE_NAME" == *.fifo ]]; then
        ALL_QUEUE_INFO+="$RESOURCE_TYPE (FIFO),$QUEUE_NAME,$REGION"$'\n'
      else
        ALL_QUEUE_INFO+="$RESOURCE_TYPE (Standard),$QUEUE_NAME,$REGION"$'\n'
      fi
    done <<< "$QUEUE_URLS"
  fi
done

echo "$ALL_QUEUE_INFO" >> "$OUTPUT_FILE"
