package lambda_queue

import (
	"aws-scalable-image-filter/internal/pkg/util"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/sqs"
)

func HandleRequest(_ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// JSON Parse Request Body
	var requestBody util.ImageDocument
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return util.InternalServerError(err), nil
	}

	// Get DynamoDB Table
	imageTable := os.Getenv("AWS_IMAGE_TABLE")
	if imageTable == "" {
		return util.InternalServerError(errors.New("Image Table was unable to be loaded from env vars.")), nil
	}

	// Get SQS Queue
	sqsQueue := os.Getenv("AWS_SQS_QUEUE")
	if sqsQueue == "" {
		return util.InternalServerError(errors.New("SQS Queue was unable to be loaded from env vars.")), nil
	}

	// Get AWS Region
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		return util.InternalServerError(errors.New("AWS Region was unable to be loaded from env vars")), nil
	}

	// Initialise AWS Session Config
	awsSessionConfig, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(awsRegion),
		},

		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return util.InternalServerError(err), nil
	}

	// Generate Document UUID and set Document Values
	documentID := uuid.New()
	requestBody.Id = documentID.String()
	requestBody.Progress = util.READY
	requestBody.DateCreated = time.Now().Unix()


	// DynamoDB Marshal Values
	imageQueueMap, err := dynamodbattribute.MarshalMap(requestBody)
	if err != nil {
		return util.InternalServerError(err), nil
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
		return util.InternalServerError(err), nil
	}

	// Start an AWS SQS Session
	awsSQSSession := sqs.New(awsSessionConfig)

	urlResult, err := awsSQSSession.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &sqsQueue,
	})
	if err != nil {
		return util.InternalServerError(err), nil
	}

	// Queue Image to SQS
	_, err = awsSQSSession.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(documentID.String()),
		QueueUrl:    urlResult.QueueUrl,
	})
	if err != nil {
		return util.InternalServerError(err), nil
	}

	// Build and return JSON response
	var queueResponse = &util.QueueResponse{
		DocumentID: documentID.String(),
	}
	response, err := json.Marshal(queueResponse)
	if err != nil {
		return util.InternalServerError(err), nil
	}

	return util.JSONStringResponse(string(response)), nil
}
