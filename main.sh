#!/bin/bash

# AWS Radar - AWS Resource Inventory Tool
# Output file name
OUTPUT_FILE="aws_resources.csv"

# Initialize the output file
> "$OUTPUT_FILE"
echo "INFO: Initialized $OUTPUT_FILE"

# Get all AWS regions
echo "INFO: Fetching AWS regions"
AWS_REGIONS=($(aws ec2 describe-regions --query "Regions[].RegionName" --output text))

# Dynamic script execution
echo "INFO: Discovering and executing AWS inventory scripts"

# Find all service directories
for SERVICE_DIR in */; do
  # Remove trailing slash
  SERVICE_NAME=${SERVICE_DIR%/}
  
  # Skip if not a directory
  if [ ! -d "$SERVICE_NAME" ]; then
    continue
  fi
  
  # Skip the directory if it contains the main script
  if [ "$SERVICE_NAME" == "$(basename "$0" .sh)" ]; then
    continue
  fi
  
  echo "INFO: Processing $SERVICE_NAME resources"
  
  # Find all script files in the directory
  SCRIPT_FILES=$(find "$SERVICE_NAME" -maxdepth 1 -name "*.sh" -type f)
  
  # If no scripts found, continue to next directory
  if [ -z "$SCRIPT_FILES" ]; then
    echo "INFO: No scripts found in $SERVICE_NAME"
    continue
  fi
  
  # Execute each script
  for SCRIPT in $SCRIPT_FILES; do
    # Make sure the script is executable
    chmod +x "$SCRIPT"
    
    # Special case for S3 which doesn't need region parameter
    if [ "$SERVICE_NAME" == "s3" ]; then
      echo "INFO: Executing $SCRIPT"
      "$SCRIPT" "$OUTPUT_FILE"
    else
      echo "INFO: Executing $SCRIPT"
      "$SCRIPT" "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
    fi
  done
done

echo "INFO: AWS Radar inventory completed successfully."
echo "INFO: Results saved to $OUTPUT_FILE"
