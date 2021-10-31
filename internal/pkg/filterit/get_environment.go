package filterit

import (
	"log"
	"os"
)

func getEnviroment() (asg string, s3Bucket string, imageTable string, sqsQueue string, redisUrl string) {
	// Get AutoScaling Group
	asg = os.Getenv("AS_GROUP")
	if asg == "" {
		log.Fatalln("Unable to find autoscaling group")
	}

	// Get S3 Bucket
	s3Bucket = os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatalln("Unable to find s3 bucket")
	}

	// Get DynamoDB Table
	imageTable = os.Getenv("AWS_IMAGE_TABLE")
	if imageTable == "" {
		log.Fatalln("Unable to find dynamodb table")
	}

	// Get SQS Queue
	sqsQueue = os.Getenv("AWS_SQS_QUEUE")
	if sqsQueue == "" {
		log.Fatalln("Unable to find sqs queue name")
	}

	// Get Redis URL
	redisUrl = os.Getenv("AWS_REDIS_ADDRESS")
	if redisUrl == "" {
		log.Fatalln("Unable to find Redis url")
	}

	return asg, s3Bucket, imageTable, sqsQueue, redisUrl
}
