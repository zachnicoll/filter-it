package lambda_feed

import (
	"aws-scalable-image-filter/internal/pkg/util"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-redis/redis/v8"
)

func HandleRequest(_ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx := context.Background()

	// Load default AWS config, including AWS_REGION env var
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return util.InternalServerError(err, "GET"), err
	}
	filter := request.QueryStringParameters["filter"]

	// If filter is not present in query, set it to the "all" redis key
	if filter == "" {
		filter = "_"
	}

	redisClient := util.ConnectToRedis()

	redisKey := filter
	cachedDoc, err := redisClient.Get(ctx, redisKey).Result()

	if err == nil && cachedDoc != "" {
		// Found a cached document for this query, return it
		return util.JSONStringResponse(cachedDoc, "GET"), nil
	} else if err != redis.Nil {
		fmt.Printf("Failed to fetch cached document from Redis: %v\n", err.Error())

		// Continue with execution, regardless of cache retrieval failure
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Scan DynamoDB table to retrieve ALL documents
	tableName := os.Getenv("AWS_IMAGE_TABLE")

	expr, err := util.BuildFilterConditions(filter)
	if err != nil {
		return util.InternalServerError(err, "GET"), err
	}

	// Perform the scan with any conditions that may be present
	scanInput := &dynamodb.ScanInput{
		TableName:        &tableName,
		FilterExpression: expr.Condition(),
	}

	fmt.Println("ScanInput: ", scanInput)

	scanOutput, err := client.Scan(ctx, scanInput)
	if err != nil {
		return util.InternalServerError(err, "GET"), err
	}

	fmt.Println("Scan Output: ", scanOutput)

	// Convert response items to list of ImageDocument
	documents := []util.ImageDocument{}

	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &documents)
	if err != nil {
		return util.InternalServerError(err, "GET"), err
	}

	// Sort documents by DateCreated, with latest first
	util.SortDocuments(documents)

	fmt.Println("Documents: ", documents)

	// Convert documents to JSON
	response, err := json.Marshal(documents)
	if err != nil {
		return util.InternalServerError(err, "GET"), err
	}

	// Convert JSON to string
	responseStr := string(response)

	util.CacheJSONString(ctx, redisKey, responseStr, redisClient)

	return util.JSONStringResponse(responseStr, "GET"), nil
}
