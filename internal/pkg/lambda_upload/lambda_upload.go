package lambda_upload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"

)

type UploadResponse struct {
	Image string `json:"image"`
	PresignURL string `json:"url"`
}

func HandleRequest() (*events.APIGatewayProxyResponse, error) {
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		fmt.Println("S3 bucket was unable to be loaded from env vars.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "S3 Bucket ENV Variable Missing.",
		}, nil
	}

	awsProfile := os.Getenv("AWS_PROFILE")
	if awsProfile == "" {
		fmt.Println("AWS Profile was unable to be loaded from env vars.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "AWS Profile ENV Variable Missing.",
		}, nil
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		fmt.Println("AWS Region was unable to be loaded from env vars.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "AWS Region ENV Variable Missing.",
		}, nil
	}

	awsSession, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: awsProfile,

		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(awsRegion),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println("Unable to configure AWS Client.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to configure AWS Client.",
		}, nil
	}

	svc := s3.New(awsSession)

	imageID := uuid.New()

	imageName := fmt.Sprintf("uploads/%s.jpg", imageID.String())
	r, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(imageName),
	})

	// Create the pre-signed url with an expiry
	url, err := r.Presign(15 * time.Minute)
	if err != nil {
		fmt.Println("Unable to generate pre-sign URL.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to generate pre-sign URL.",
		}, nil
	}

	var uploadReponse = &UploadResponse{
		Image: imageID.String(),
		PresignURL: url,
	}

	resp, err := json.Marshal(uploadReponse)
	if err != nil {
		fmt.Println("Unable to convert response to JSON.")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to convert response to JSON.",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string {
			"Content-Type": "application/json",
		},
		Body: string(resp),
	}, nil
}