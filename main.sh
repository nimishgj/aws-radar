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

# Call ec2_elastic_ips.sh with output file and regions
./ec2_elastic_ips.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ec2_key_pairs.sh with output file and regions
./ec2_key_pairs.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ec2_network_interfaces.sh with output file and regions
./ec2_network_interfaces.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ec2_load_balancers.sh with output file and regions
./ec2_load_balancers.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ec2_target_groups.sh with output file and regions
./ec2_target_groups.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ec2_asg.sh with output file and regions
./ec2_asg.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc.sh with output file and regions
./vpc.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_subnets.sh with output file and regions
./vpc_subnets.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_route_table.sh with output file and regions
./vpc_route_table.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_internet_gateway.sh with output file and regions
./vpc_internet_gateway.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_egress_igw.sh with output file and regions
./vpc_egress_igw.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_dhpc_option_set.sh with output file and regions
./vpc_dhpc_option_set.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_managed_prefix_list.sh with output file and regions
./vpc_managed_prefix_list.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_nat_gateway.sh with output file and regions
./vpc_nat_gateway.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_network_acls.sh with output file and regions
./vpc_network_acls.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_customer_gateway.sh with output file and regions
./vpc_customer_gateway.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_virtual_private_gateway.sh with output file and regions
./vpc_virtual_private_gateway.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_site_to_site_vpn.sh with output file and regions
./vpc_site_to_site_vpn.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call vpc_transit_gateways.sh with output file and regions
./vpc_transit_gateways.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ecs_cluster.sh with output file and regions
./ecs_cluster.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ecs_namespaces.sh with output file and regions
./ecs_namespaces.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ecs_task_definitions.sh with output file and regions
./ecs_task_definitions.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ecr_public_repositories.sh with output file and regions
./ecr_public_repositories.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"

# Call ecr_private_repositories.sh with output file and regions
./ecr_private_repositories.sh "$OUTPUT_FILE" "${AWS_REGIONS[@]}"
echo "INFO: Done."
