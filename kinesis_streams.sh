#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Kinesis Streams"
ALL_STREAM_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Kinesis Streams in $REGION"
  STREAMS=$(aws kinesis list-streams --region "$REGION" \
    --query "StreamNames[]" \
    --output text 2>/dev/null || echo "")

  for STREAM_NAME in $STREAMS; do
    ALL_STREAM_INFO+="$RESOURCE_TYPE,$STREAM_NAME,$REGION"$'\n'
  done
done

echo "$ALL_STREAM_INFO" >> "$OUTPUT_FILE"
