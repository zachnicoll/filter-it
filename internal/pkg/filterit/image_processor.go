package filterit

import (
	"aws-scalable-image-filter/internal/pkg/util"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

func applyFilters(filters []int) {
	// TODO: Use image magik to apply each of the filters to the image
}

func processMessage(id string) {
	// TODO: Fetch document form DynamoDb with respective id

	// TODO: Set document's progress attribute to PROCESSING

	// TODO: Fetch image from S3 based on the document's image attribute

	// TODO: Apply filters to image

	// TODO: Re-upload filtered image to S3

	// TODO: Write new filenname to document's image attribute

	// TODO: Set document progress attribute to DONE

	// TODO: Invalid cache with all keys containing filters (use KEYS command)
}

func WatchQueue() {
	ctx := context.Background()

	// Load default AWS config, including AWS_REGION env var
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		util.FatalLog("Failed to load default AWS config!", err)
	}

	// Get instance ID
	instanceID := util.FetchInstanceID()

	asg, s3Bucket, imageTable, sqsQueue := getEnviroment()

	// AWS DynamoDB session
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// AWS SQS session
	sqsClient := sqs.NewFromConfig(cfg)

	// AWS S3 session
	s3Client := s3.NewFromConfig(cfg)

	// AWS ASG session
	asgClient := autoscaling.NewFromConfig(cfg)

	// SQS URL
	urlResult, err := sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &sqsQueue,
	})
	if err != nil {
		util.FatalLog("Unable to fetch aws sqs url", err)
	}

	for {
		// SQS Intake Parameters
		params := &sqs.ReceiveMessageInput{
			QueueUrl:            urlResult.QueueUrl, // Required
			MaxNumberOfMessages: 1,
			MessageAttributeNames: []string{
				"All",
			},
			WaitTimeSeconds: 30,
		}

		// Receive an SQS Message
		resp, err := sqsClient.ReceiveMessage(ctx, params)
		if err != nil {
			util.FatalLog("Unable to fetch aws sqs message", err)
		}

		// Check a message was received
		if len(resp.Messages) == 1 {
			_, err = asgClient.SetInstanceProtection(ctx, &autoscaling.SetInstanceProtectionInput{
				InstanceIds:          []string{instanceID},
				AutoScalingGroupName: aws.String(asg),
				ProtectedFromScaleIn: aws.Bool(true),
			})
			if err != nil {
				util.FatalLog("Unable to enable scale-in protection", err)
			}

			targetMessage := resp.Messages[0]

			_, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      urlResult.QueueUrl,          // Required
				ReceiptHandle: targetMessage.ReceiptHandle, // Required
			})
			if err != nil {
				util.FatalLog("Unable to delete aws sqs message", err)
			}

			// Get DynamoDB image info from message body
			result, err := dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
				TableName: aws.String(imageTable),
				Key: map[string]types.AttributeValue{
					"id": types.AttributeValueMemberS{
						Value: *targetMessage.Body,
					},
					"progress": types.AttributeValueMemberN{
						Value: string(util.READY),
					},
				},
			})
			if err != nil {
				util.SafeFail(sqsClient, dynamoClient, asgClient, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
				util.FatalLog("Unable to check aws sqs dynamodb status", err)
			}

			// Unmarshal image document
			var imageDocument util.ImageDocument
			err = attributevalue.UnmarshalMap(result.Item, &imageDocument)
			if err != nil {
				util.SafeFail(sqsClient, dynamoClient, asgClient, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
				util.FatalLog("Unable to unmarshal aws sqs dynamodb status", err)
			}

			imageDocument.Progress = util.PROCESSING

			imageDocumentMap, err := attributevalue.MarshalMap(imageDocument)
			if err != nil {
				util.SafeFail(sqsClient, dynamoClient, asgClient, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
				util.FatalLog("Unable to marshal aws sqs dynamodb (processing)", err)
			}

			_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
				Item:      imageDocumentMap,
				TableName: aws.String(imageTable),
			})
			if err != nil {
				util.SafeFail(sqsClient, dynamoClient, asgClient, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
				util.FatalLog("Unable to update aws sqs dynamodb (processing)", err)
			}

			// TODO: IMAGE
			_, err = s3Client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(s3Bucket),
				Key:    aws.String(imageDocument.Image),
			})
			if err != nil {
				util.SafeFail(sqsClient, dynamoClient, asgClient, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
				util.FatalLog("Unable to get S3 image", err)
			}

			// PROCESS IMAGE

			// Generate Image UUID
			imageID := uuid.New()

			// String format S3 image name and generate a S3 put object
			imageName := fmt.Sprintf("uploads/%s.jpg", imageID.String())
			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(s3Bucket),
				Key:    aws.String(imageName),
				// TODO: IMAGE
				//Body: ,
			})
			if err != nil {
				util.SafeFail(sqsClient, dynamoClient, asgClient, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
				util.FatalLog("Unable to put new S3 image", err)
			}

			imageDocument.Progress = util.DONE
			imageDocument.Image = imageName

			imageDocumentMap, err = attributevalue.MarshalMap(imageDocument)
			if err != nil {
				util.SafeFail(sqsClient, dynamoClient, asgClient, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
				util.FatalLog("Unable to marshal aws sqs dynamodb (processed)", err)
			}

			_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
				Item:      imageDocumentMap,
				TableName: aws.String(imageTable),
			})
			if err != nil {
				util.SafeFail(sqsClient, dynamoClient, asgClient, asg, instanceID, imageTable, *urlResult.QueueUrl, *targetMessage.Body)
				util.FatalLog("Unable to update aws sqs dynamodb (processed)", err)
			}
		}
	}

	// TODO: In a loop, check if the SQS queue has a new message

	// TODO: If message, spin off a subroutine and process the message - processMessage(id)

	// TODO: Log custom metric to CloudWatch

}
