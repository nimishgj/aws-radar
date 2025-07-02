#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Load Balancers"
ALL_LB_INFO=""

for REGION in "${REGIONS[@]}"; do
  echo "INFO: Fetching Load Balancers in $REGION"
  # Get classic load balancers
  CLASSIC_LBS=$(aws elb describe-load-balancers --region "$REGION" \
    --query "LoadBalancerDescriptions[].LoadBalancerName" \
    --output text 2>/dev/null || echo "")
  
  # Get application and network load balancers
  ALB_NLB=$(aws elbv2 describe-load-balancers --region "$REGION" \
    --query "LoadBalancers[].LoadBalancerName" \
    --output text 2>/dev/null || echo "")
  
  # Process classic load balancers
  for NAME in $CLASSIC_LBS; do
    ALL_LB_INFO+="$RESOURCE_TYPE-Classic,$NAME,$REGION"$'\n'
  done
  
  # Process application and network load balancers
  for NAME in $ALB_NLB; do
    ALL_LB_INFO+="$RESOURCE_TYPE,$NAME,$REGION"$'\n'
  done
done

echo "$ALL_LB_INFO" >> "$OUTPUT_FILE"
