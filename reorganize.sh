#!/bin/bash

# Create directories for each service
mkdir -p ec2 vpc cloudwatch ecs ecr rds lambda elasticache msk apigateway s3 eks sns sqs secretsmanager kms ssm route53 dynamodb kinesis stepfunctions amazon_mq

# Move files to their respective directories
mv ec2_*.sh ec2/
mv vpc_*.sh vpc/
mv cloudwatch_*.sh cloudwatch/
mv ecs_*.sh ecs/
mv ecr_*.sh ecr/
mv rds*.sh rds/
mv lambda_*.sh lambda/
mv elasticache_*.sh elasticache/
mv msk_*.sh msk/
mv apigateway_*.sh apigateway/
mv s3.sh s3/
mv eks_*.sh eks/
mv sns_*.sh sns/
mv sqs_*.sh sqs/
mv secretsmanager_*.sh secretsmanager/
mv kms_*.sh kms/
mv ssm_*.sh ssm/
mv route53_*.sh route53/
mv dynamodb_*.sh dynamodb/
mv kinesis_*.sh kinesis/
mv stepfunctions_*.sh stepfunctions/
mv amazon_mq_*.sh amazon_mq/

# Rename the files to remove service prefix (optional)
for dir in */; do
  cd $dir
  for file in *.sh; do
    # Extract the part after the service name
    new_name=$(echo $file | sed 's/^[^_]*_//')
    # Only rename if the pattern matched
    if [[ "$new_name" != "$file" ]]; then
      mv "$file" "$new_name"
    fi
  done
  cd ..
done

echo "Files have been reorganized into service directories"
