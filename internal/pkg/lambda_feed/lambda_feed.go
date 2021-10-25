package lambda_feed

import (
	"aws-scalable-image-filter/internal/pkg/helpers"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sort"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type FeedRequestBody struct {
	Filters []string `json:"filters"`
}

const (
	GREYSCALE int = 0
	SEPIA         = 1
	INVERT        = 2
)

const (
	READY      int = 0
	PROCESSING     = 1
	DONE           = 2
	FAILED         = 3
)

type ImageDocument struct {
	Id          string `json:"id"`
	DateCreated int    `json:"date_created"`
	Filters     []int  `json:"filters"`
	Progress    int    `json:"progress"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Image       string `json:"image"`
}

func HandleRequest(ctx context.Context, event FeedRequestBody) *events.APIGatewayProxyResponse {
	// Load default AWS config, including AWS_REGION env var
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return helpers.InternalServerError(err)
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Scan DynamoDB table to retrieve ALL documents
	tableName := os.Getenv("AWS_IMAGE_TABLE")

	input := &dynamodb.ScanInput{
		TableName: &tableName,
	}

	resp, err := client.Scan(context.TODO(), input)
	if err != nil {
		return helpers.InternalServerError(err)
	}

	// Convert response items to list of ImageDocument
	documents := []ImageDocument{}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &documents)
	if err != nil {
		return helpers.InternalServerError(err)
	}

	// Sort documents by DateCreated, with latest first
	sort.Slice(documents, func(i, j int) bool {
		return documents[i].DateCreated > documents[j].DateCreated
	})

	// Convert documents to JSON
	response, err := json.Marshal(documents)
	if err != nil {
		return helpers.InternalServerError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}
}
