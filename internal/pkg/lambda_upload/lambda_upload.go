package lambda_upload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
)

type UploadResponse struct {
	Image      string `json:"image"`
	PresignURL string `json:"url"`
}

func HandleRequest() (*events.APIGatewayProxyResponse, error) {
	// Get S3 Bucket Name
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		fmt.Println("S3 bucket was unable to be loaded from env vars.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "S3 Bucket ENV Variable Missing.",
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

	// Start an AWS Session
	awsSession := s3.New(awsSessionConfig)

	// Generate Image UUID
	imageID := uuid.New()

	// String format S3 image name and generate a S3 put request
	imageName := fmt.Sprintf("uploads/%s.jpg", imageID.String())
	s3PutRequest, _ := awsSession.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(imageName),
	})

	// Fetch S3 pre-sign URL from put request
	presignURL, err := s3PutRequest.Presign(15 * time.Minute)
	if err != nil {
		fmt.Println("Unable to generate pre-sign URL.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to generate pre-sign URL.",
		}, nil
	}

	// URL escape format
	cleanURL := url.QueryEscape(presignURL)

	// Build and return JSON response
	var uploadResponse = &UploadResponse{
		Image:      imageName,
		PresignURL: cleanURL,
	}
	response, err := json.Marshal(uploadResponse)
	if err != nil {
		fmt.Println("Unable to convert response to JSON.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to convert response to JSON.",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil
}
