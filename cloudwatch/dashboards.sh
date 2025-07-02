#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch Dashboards"
ALL_DASHBOARD_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch Dashboards in $REGION"
  DASHBOARDS=$(aws cloudwatch list-dashboards --region "$REGION" \
    --query "DashboardEntries[].DashboardName" \
    --output text 2>/dev/null || echo "")

  for NAME in $DASHBOARDS; do
    ALL_DASHBOARD_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_DASHBOARD_INFO" >> "$OUTPUT_FILE"
