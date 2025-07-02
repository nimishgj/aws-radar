#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch RUM Apps"
ALL_RUM_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch RUM App Monitors in $REGION"
  APP_MONITORS=$(aws rum list-app-monitors --region "$REGION" \
    --query "AppMonitorSummaries[].Name" \
    --output text 2>/dev/null || echo "")

  for NAME in $APP_MONITORS; do
    ALL_RUM_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_RUM_INFO" >> "$OUTPUT_FILE"
