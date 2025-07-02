#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Amazon MQ Brokers"
ALL_BROKER_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Amazon MQ Brokers in $REGION"
  BROKERS=$(aws mq list-brokers --region "$REGION" \
    --query "BrokerSummaries[].BrokerId" \
    --output text 2>/dev/null || echo "")

  for ID in $BROKERS; do
    # Get the broker name from the broker ID
    BROKER_NAME=$(aws mq describe-broker --broker-id "$ID" --region "$REGION" \
      --query "BrokerName" --output text 2>/dev/null || echo "")
    
    # If we got the broker name, use it; otherwise use the ID
    if [ ! -z "$BROKER_NAME" ] && [ "$BROKER_NAME" != "None" ]; then
      ALL_BROKER_INFO+="$RESOURCE_TYPE,$BROKER_NAME,$REGION"$'\n'
    else
      ALL_BROKER_INFO+="$RESOURCE_TYPE,$ID,$REGION"$'\n'
    fi
  done
done

echo "$ALL_BROKER_INFO" >> "$OUTPUT_FILE"
