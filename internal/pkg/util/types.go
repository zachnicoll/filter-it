package util

import (
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-redis/redis/v8"
)

const (
	GREYSCALE int = 0
	SEPIA         = 1
	INVERT        = 2
)

const (
	READY      int = 0
	PROCESSING     = 1
	DONE           = 2
	FAILED         = 3
)

type ImageDocument struct {
	Id          string `json:"id,omitempty"`
	DateCreated int64  `json:"date_created,omitempty"`
	Tag         int    `json:"tag"`
	Progress    int    `json:"progress,omitempty"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Image       string `json:"image"`
}

type QueueResponse struct {
	DocumentID string `json:"id"`
}

type Clients struct {
	DynamoDb *dynamodb.Client
	SQS      *sqs.Client
	ASG      *autoscaling.Client
	S3       *s3.Client
	Redis    *redis.Client
}

type MetaData struct {
	ImageTable *string
	InstanceID *string
	SQSUrl     *string
	S3Bucket   *string
	ASGName    *string
}
