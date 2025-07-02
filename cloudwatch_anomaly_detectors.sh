#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="CloudWatch Anomaly Detectors"
ALL_DETECTOR_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching CloudWatch Anomaly Detectors in $REGION"
  # For anomaly detectors we need to capture namespace/metric combinations
  DETECTORS=$(aws cloudwatch describe-anomaly-detectors --region "$REGION" \
    --query "AnomalyDetectors[].{namespace:Namespace,metric:MetricName}" \
    --output text 2>/dev/null || echo "")

  if [ ! -z "$DETECTORS" ]; then
    while read -r NAMESPACE METRIC; do
      DETECTOR_NAME="${NAMESPACE}/${METRIC}"
      ALL_DETECTOR_INFO+="$RESOURCE_TYPE,$DETECTOR_NAME,$REGION"$'\n'
    done <<< "$DETECTORS"
  fi
done

echo "$ALL_DETECTOR_INFO" >> "$OUTPUT_FILE"
