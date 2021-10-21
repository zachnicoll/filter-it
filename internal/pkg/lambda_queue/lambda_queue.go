package lambda_queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type QueueRequest struct {
	Id string `json:"id,omitempty"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Filters []string `json:"filter"`
	Image string `json:"image"`
	Progress string `json:"progress,omitempty"`
}

type QueueResponse struct {
	DocumentID string `json:"id"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// JSON Parse Request Body
	var requestBody QueueRequest
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		fmt.Println("Bad Request")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Malformed Body",
		}, nil
	}

	// Get DynamoDB Table
	imageTable := os.Getenv("AWS_IMAGE_TABLE")
	if imageTable == "" {
		fmt.Println("Image Table was unable to be loaded from env vars.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Image Table ENV Variable Missing.",
		}, nil
	}

	// Get SQS Queue
	sqsQueue := os.Getenv("AWS_SQS_QUEUE")
	if sqsQueue == "" {
		fmt.Println("SQS Queue was unable to be loaded from env vars.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "SQS Queue ENV Variable Missing.",
		}, nil
	}

	// Get AWS Region
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		fmt.Println("AWS Region was unable to be loaded from env vars.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "AWS Region ENV Variable Missing.",
		}, nil
	}

	// Initialise AWS Session Config
	awsSessionConfig, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(awsRegion),
		},

		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println("Unable to configure AWS Client.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to configure AWS Client.",
		}, nil
	}

	// Generate Document UUID and set Document Values
	documentID := uuid.New()
	requestBody.Id = documentID.String()
	requestBody.Progress = "PROCESSING"

	// DynamoDB Marshal Values
	imageQueueMap, err := dynamodbattribute.MarshalMap(requestBody)
	if err != nil {
		fmt.Println("Unable to marshal dynamodb fields.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to marshal dynamodb fields.",
		}, nil
	}

	// Start an AWS DB Session
	awsDBSession := dynamodb.New(awsSessionConfig)

	// Initialise DynamoDB PutItem Config
	input := &dynamodb.PutItemInput{
		Item:      imageQueueMap,
		TableName: aws.String(imageTable),
	}

	// Put DynamoDB Item
	_, err = awsDBSession.PutItem(input)
	if err != nil {
		fmt.Println("Unable to write to dynamodb.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to write to dynamodb.",
		}, nil
	}

	// Start an AWS SQS Session
	awsSQSSession := sqs.New(awsSessionConfig)

	// Queue Image to SQS
	_, err = awsSQSSession.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(documentID.String()),
		QueueUrl:    &sqsQueue,
	})
	if err != nil {
		fmt.Println("Unable to queue message to SQS.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to queue message to SQS.",
		}, nil
	}

	// Build and return JSON response
	var queueResponse = &QueueResponse{
		DocumentID: documentID.String(),
	}
	response, err := json.Marshal(queueResponse)
	if err != nil {
		fmt.Println("Unable to convert response to JSON.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to convert response to JSON.",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: string(response),
	}, nil
}
