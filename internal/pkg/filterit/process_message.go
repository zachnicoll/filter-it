package filterit

import (
	"aws-scalable-image-filter/internal/pkg/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

func processMessage(wg *sync.WaitGroup, ctx context.Context, msg *sqsTypes.Message, clients *util.Clients, metaData *util.MetaData) {
	defer wg.Done()

	log.Println("Processing Queue Message", *msg.Body)

	var queueMsg util.QueueResponse
	jsonStr := *msg.Body
	err := json.Unmarshal([]byte(jsonStr), &queueMsg)

	if err != nil {
		util.SafeFailAndLog(clients, metaData, &queueMsg, "Could not unmarshal queue message", err)
	}

	log.Println("Getting image document from DynamoDb...")

	// Get DynamoDB image info from message body
	result, err := clients.DynamoDb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: metaData.ImageTable,
		Key: map[string]dynamoTypes.AttributeValue{
			"id": &dynamoTypes.AttributeValueMemberS{
				Value: queueMsg.DocumentID,
			},
			"date_created": &dynamoTypes.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", queueMsg.DateCreated),
			},
		},
	})
	if err != nil {
		util.SafeFailAndLog(clients, metaData, &queueMsg, "Unable to check aws sqs dynamodb status", err)
	}

	// Unmarshal image document
	var imageDocument util.ImageDocument
	err = attributevalue.UnmarshalMap(result.Item, &imageDocument)
	if err != nil {
		util.SafeFailAndLog(clients, metaData, &queueMsg, "Unable to unmarshal aws sqs dynamodb status", err)
	}

	log.Println("Updating image document in DynamoDb with Processing progress...")

	// Mark image document as processing
	imageDocument.Progress = util.PROCESSING

	err = util.UpdateDocument(ctx, clients, metaData.ImageTable, &imageDocument)
	if err != nil {
		util.SafeFailAndLog(clients, metaData, &queueMsg, "Unable to update aws sqs dynamodb (processing)", err)
	}

	log.Println("Getting original image from S3...")

	// Get image from S3
	s3Object, err := clients.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: metaData.S3Bucket,
		Key:    aws.String(imageDocument.Image),
	})
	if err != nil {
		util.SafeFailAndLog(clients, metaData, &queueMsg, "Unable to get S3 image", err)
	}

	// Apply filter to image
	log.Println("Applying Filter To Image:", imageDocument.Id, "Filter:", imageDocument.Tag)

	blob, err := applyFilter(s3Object.Body, imageDocument.Tag)
	if err != nil {
		util.SafeFailAndLog(clients, metaData, &queueMsg, "Failed to filter image", err)
	}

	log.Println("Successfully Applied Filter", imageDocument.Id)

	reader := bytes.NewReader(blob)

	// New image as UUID
	imageName := fmt.Sprintf("uploads/%s.jpg", uuid.New())

	// Upload new filtered image
	_, err = clients.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: metaData.S3Bucket,
		Key:    aws.String(imageName),
		Body:   reader,
	})
	if err != nil {
		util.SafeFailAndLog(clients, metaData, &queueMsg, "Unable to put new S3 image", err)
	}

	log.Println("Uploaded Image to S3", imageName)

	// Mark image as DONE processing
	imageDocument.Progress = util.DONE

	// Assign new image name to point to filtered image
	imageDocument.Image = imageName

	err = util.UpdateDocument(ctx, clients, metaData.ImageTable, &imageDocument)
	if err != nil {
		util.SafeFailAndLog(clients, metaData, &queueMsg, "Unable to update aws sqs dynamodb (processing)", err)
	}

	log.Println("Finished Processing Image", imageDocument.Id)
}
