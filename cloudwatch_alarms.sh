#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch Alarms"
ALL_ALARM_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch Alarms in $REGION"
  ALARMS=$(aws cloudwatch describe-alarms --region "$REGION" \
    --query "MetricAlarms[].AlarmName" \
    --output text 2>/dev/null || echo "")

  for NAME in $ALARMS; do
    ALL_ALARM_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_ALARM_INFO" >> "$OUTPUT_FILE"
