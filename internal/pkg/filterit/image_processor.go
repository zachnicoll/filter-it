package filterit

import (
	"aws-scalable-image-filter/internal/pkg/util"

	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func WatchQueue() {
	ctx := context.Background()
	var wg sync.WaitGroup

	var instanceProtection = false

	// Load default AWS config, including AWS_REGION env var
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		util.FatalLog("Failed to load default AWS config!", err)
	}

	// Get instance ID
	instanceID := "" //util.FetchInstanceID()

	asg, s3Bucket, imageTable, sqsQueue := getEnvironment()

	dynamoClient, sqsClient, s3Client, asgClient, redisClient := getClients(cfg)

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
		Redis:    redisClient,
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
		WaitTimeSeconds: 10,
	}

	// Intake forever
	for {
		// Max of 2 messages per intake
		for i := 0; i < 2; i++ {
			// Receive an SQS Message
			resp, err := sqsClient.ReceiveMessage(ctx, receiveMessageInput)
			if err != nil {
				util.FatalLog("Unable to fetch aws sqs message", err)
			}

			// Check that message was received
			if len(resp.Messages) == 1 {
				if !instanceProtection {
					// Protect this instance from being destroyed while processing
					// _, err = asgClient.SetInstanceProtection(ctx, &autoscaling.SetInstanceProtectionInput{
					// 	InstanceIds:          []string{instanceID},
					// 	AutoScalingGroupName: aws.String(asg),
					// 	ProtectedFromScaleIn: aws.Bool(true),
					// })
					// if err != nil {
					// 	util.FatalLog("Unable to enable scale-in protection", err)
					// }
					// instanceProtection = true
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

				// Offload to a coroutine
				wg.Add(1)
				go processMessage(&wg, ctx, &targetMessage, &clients, &metaData)
			}
		}

		wg.Wait()

		// Unprotect this instance from being destroyed while processing
		// _, err = asgClient.SetInstanceProtection(ctx, &autoscaling.SetInstanceProtectionInput{
		// 	InstanceIds:          []string{instanceID},
		// 	AutoScalingGroupName: aws.String(asg),
		// 	ProtectedFromScaleIn: aws.Bool(false),
		// })
		// if err != nil {
		// 	util.FatalLog("Unable to enable scale-in protection", err)
		// }
		// instanceProtection = true
	}
}
