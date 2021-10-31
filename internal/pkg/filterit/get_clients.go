package filterit

import (
	"aws-scalable-image-filter/internal/pkg/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-redis/redis/v8"
)

func getClients(cfg aws.Config) (DDB *dynamodb.Client, SQS *sqs.Client, S3 *s3.Client, ASG *autoscaling.Client, Redis *redis.Client) {
	DDB = dynamodb.NewFromConfig(cfg)
	SQS = sqs.NewFromConfig(cfg)
	S3 = s3.NewFromConfig(cfg)
	ASG = autoscaling.NewFromConfig(cfg)
	Redis = util.ConnectToRedis()

	return DDB, SQS, S3, ASG, Redis
}
