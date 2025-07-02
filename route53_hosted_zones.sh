#!/bin/bash

OUTPUT_FILE="$1"
shift
REGIONS=("$@")

RESOURCE_TYPE="Route53 Hosted Zones"
ALL_ZONE_INFO=""

# Route53 is a global service but we'll include it for consistency
echo "INFO: Fetching Route53 Hosted Zones"
ZONES=$(aws route53 list-hosted-zones \
  --query "HostedZones[].[Id,Name]" \
  --output text 2>/dev/null || echo "")

if [ -n "$ZONES" ]; then
  while read -r ZONE_ID ZONE_NAME; do
    # Remove trailing dot from zone name
    ZONE_NAME=${ZONE_NAME%%.}
    # Clean up zone ID to remove /hostedzone/ prefix
    ZONE_ID=${ZONE_ID##*/}
    ALL_ZONE_INFO+="$RESOURCE_TYPE,$ZONE_NAME ($ZONE_ID),global"$'\n'
  done <<< "$ZONES"
fi

echo "$ALL_ZONE_INFO" >> "$OUTPUT_FILE"
