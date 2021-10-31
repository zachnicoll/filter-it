#!/usr/bin/env sh

echo "S3: $S3_BUCKET"
echo "DDB: $AWS_IMAGE_TABLE"
echo "SQS: $AWS_SQS_QUEUE"
echo "REDIS: $AWS_REDIS_ADDRESS"
echo "AKID: $AWS_ACCESS_KEY_ID"
echo "AK: $AWS_SECRET_ACCESS_KEY"
echo "TOKEN: $AWS_SESSION_TOKEN"

./filterit