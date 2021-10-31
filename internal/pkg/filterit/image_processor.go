package filterit

import (
	"aws-scalable-image-filter/internal/pkg/util"
	"context"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func WatchQueue() {
	ctx := context.Background()
	runtime.GOMAXPROCS(2)

	// Load default AWS config, including AWS_REGION env var
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		util.FatalLog("Failed to load default AWS config!", err)
	}

	// Get instance ID
	instanceID := util.FetchInstanceID()

	asg, s3Bucket, imageTable, sqsQueue, _ := getEnvironment()

	dynamoClient, sqsClient, s3Client, asgClient := getClients(cfg)

	// SQS URL
	urlResult, err := sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &sqsQueue,
	})
	if err != nil {
		util.FatalLog("Unable to fetch aws sqs url", err)
	}

	clients := util.Clients{
		DynamoDb: dynamoClient,
		SQS:      sqsClient,
		S3:       s3Client,
		ASG:      asgClient,
	}

	metaData := util.MetaData{
		S3Bucket:   &s3Bucket,
		ImageTable: &imageTable,
		SQSUrl:     urlResult.QueueUrl,
		InstanceID: &instanceID,
		ASGName:    &asg,
	}

	// SQS Intake Parameters
	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl:            urlResult.QueueUrl, // Required
		MaxNumberOfMessages: 1,
		MessageAttributeNames: []string{
			"All",
		},
		WaitTimeSeconds: 30,
	}

	for {
		// Receive an SQS Message
		resp, err := sqsClient.ReceiveMessage(ctx, receiveMessageInput)
		if err != nil {
			util.FatalLog("Unable to fetch aws sqs message", err)
		}

		// Check a message was received
		if len(resp.Messages) == 1 {
			// Protect this instance from being destroyed while processing
			_, err = asgClient.SetInstanceProtection(ctx, &autoscaling.SetInstanceProtectionInput{
				InstanceIds:          []string{instanceID},
				AutoScalingGroupName: aws.String(asg),
				ProtectedFromScaleIn: aws.Bool(true),
			})
			if err != nil {
				util.FatalLog("Unable to enable scale-in protection", err)
			}

			targetMessage := resp.Messages[0]

			// Remove message from queue now that we have it
			_, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      urlResult.QueueUrl,          // Required
				ReceiptHandle: targetMessage.ReceiptHandle, // Required
			})
			if err != nil {
				util.FatalLog("Unable to delete aws sqs message", err)
			}

			go processMessage(ctx, &targetMessage, &clients, &metaData)
		}
	}
}
