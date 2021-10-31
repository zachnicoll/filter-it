package filterit

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func getClients(cfg aws.Config) (DDB *dynamodb.Client, SQS *sqs.Client, S3 *s3.Client, ASG *autoscaling.Client) {
	DDB = dynamodb.NewFromConfig(cfg)
	SQS = sqs.NewFromConfig(cfg)
	S3 = s3.NewFromConfig(cfg)
	ASG = autoscaling.NewFromConfig(cfg)

	return DDB, SQS, S3, ASG
}
