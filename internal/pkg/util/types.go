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
	Id          string `json:"id,omitempty" dynamodbav:"id"`
	DateCreated int64  `json:"date_created,omitempty" dynamodbav:"date_created"`
	Tag         int    `json:"tag" dynamodbav:"tag"`
	Progress    int    `json:"progress,omitempty" dynamodbav:"progress"`
	Title       string `json:"title" dynamodbav:"title"`
	Author      string `json:"author" dynamodbav:"author"`
	Image       string `json:"image" dynamodbav:"image"`
	ImageURL    string `json:"image_url,omitempty"`
}

type QueueResponse struct {
	DocumentID  string `json:"id"`
	DateCreated int64  `json:"date_created,omitempty"`
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

type ProgressRequest struct {
	DocumentID string `json:"id"`
}

type ProgressResponse struct {
	DocumentID string `json:"id"`
	ImageURL   string `json:"imageurl"`
	Progress   int    `json:"progress"`
}
