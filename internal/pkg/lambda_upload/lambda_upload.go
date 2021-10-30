package lambda_upload

import (
	"aws-scalable-image-filter/internal/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
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
		return util.InternalServerError(errors.New("S3 Bucket was unable to be loaded from env vars")), nil
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
		return util.InternalServerError(err), nil
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
		return util.InternalServerError(err), nil
	}

	return util.JSONStringResponse(string(response)), nil
}
