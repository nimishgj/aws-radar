#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch Internet Monitor"
ALL_MONITOR_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch Internet Monitors in $REGION"
  MONITORS=$(aws internetmonitor list-monitors --region "$REGION" \
    --query "Monitors[].MonitorName" \
    --output text 2>/dev/null || echo "")

  for NAME in $MONITORS; do
    ALL_MONITOR_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_MONITOR_INFO" >> "$OUTPUT_FILE"
