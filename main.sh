#!/bin/bash

# Output file name
OUTPUT_FILE="aws_resources.csv"

# Initialize the output file
> "$OUTPUT_FILE"
echo "INFO: Initialized $OUTPUT_FILE"

# Get all AWS regions
echo "INFO: Fetching AWS regions"
AWS_REGIONS=($(aws ec2 describe-regions --query "Regions[].RegionName" --output text))

# Pass regions and output file to ec2.sh
./ec2_instances.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call s3.sh with output file
./s3.sh "$OUTPUT_FILE"

# Call ec2_ebs_volumes.sh with output file and regions
./ec2_ebs_volumes.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ec2_ebs_snapshots.sh with output file and regions
./ec2_ebs_snapshots.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ec2_security_group.sh with output file and regions
./ec2_security_group.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
echo "INFO: Done."
