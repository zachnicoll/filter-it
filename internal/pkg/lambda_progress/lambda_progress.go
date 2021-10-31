package lambda_progress

import (
	"aws-scalable-image-filter/internal/pkg/util"
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-lambda-go/events"
)

func HandleRequest(_ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var requestBody util.ProgressRequest
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return util.InternalServerError(err), nil
	}

	// Get S3 Bucket Name
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		return util.InternalServerError(errors.New("S3 Bucket was unable to be loaded from env vars")), nil
	}

	// Get DynamoDB Table
	imageTable := os.Getenv("AWS_IMAGE_TABLE")
	if imageTable == "" {
		return util.InternalServerError(errors.New("Image Table was unable to be loaded from env vars.")), nil
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

	// Start an AWS DB Session
	awsDBSession := dynamodb.New(awsSessionConfig)

	// Start an AWS s3 Session
	awsS3Session := s3.New(awsSessionConfig)

	// Get DynamoDB image info from message body
	result, err := awsDBSession.Query(&dynamodb.QueryInput{
		TableName: aws.String(imageTable),
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(requestBody.DocumentID),
					},
				},
			},
		},
	})
	if err != nil {
		return util.InternalServerError(err), nil
	}

	// Unmarshal image document
	var imageDocument []util.ImageDocument
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &imageDocument)
	if err != nil {
		return util.InternalServerError(err), nil
	}

	if len(imageDocument) == 0 {
		return util.InternalServerError(errors.New("No matching image document found.")), nil
	} else if len(imageDocument) > 1 {
		return util.InternalServerError(errors.New("Returned multiple image documents for that ID.")), nil
	}

	progressResponse := &util.ProgressResponse{
		DocumentID: imageDocument[0].Id,
		Progress: imageDocument[0].Progress,
	}

	if imageDocument[0].Progress == util.DONE {
		s3GetRequest, _ := awsS3Session.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(imageDocument[0].Image),
		})

		presignURL, err := s3GetRequest.Presign(24 * time.Hour)
		if err != nil {
			return util.InternalServerError(err), nil
		}

		// URL escape format
		cleanURL := url.QueryEscape(presignURL)

		progressResponse.ImageURL = cleanURL
	}

	response, err := json.Marshal(progressResponse)
	if err != nil {
		return util.InternalServerError(err), nil
	}

	return util.JSONStringResponse(string(response)), nil
}
