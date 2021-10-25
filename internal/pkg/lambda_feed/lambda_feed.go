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
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

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

	// Create a condition for each filter supplied with the request,
	// so that each retrieved document must include each of the supplied filters
	conditions := []expression.ConditionBuilder{}

	for _, filter := range event.Filters {
		condition := expression.Name("filters").In(expression.Value(filter))
		conditions = append(conditions, condition)
	}

	builder := expression.NewBuilder()

	for _, condition := range conditions {
		builder = builder.WithCondition(condition)
	}

	expr, err := builder.Build()
	if err != nil {
		return helpers.InternalServerError(err)
	}

	// Perform the scan with any conditions that may be present
	input := &dynamodb.ScanInput{
		TableName:        &tableName,
		FilterExpression: expr.Condition(),
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
