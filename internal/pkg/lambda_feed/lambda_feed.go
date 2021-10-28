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

func HandleRequest(ctx context.Context, event util.FeedRequestBody) (*events.APIGatewayProxyResponse, error) {
	// Load default AWS config, including AWS_REGION env var
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return util.InternalServerError(err), err
	}

	filters := event.Filters
	util.SortFilters(filters)

	redisClient := util.ConnectToRedis()

	// Attempt to get cached documents from Redis
	redisKey := util.ConstructRedisKey(filters)
	cachedDoc, err := redisClient.Get(ctx, redisKey).Result()

	if err == nil && cachedDoc != "" {
		// Found a cached document for this query, return it
		return util.JSONStringResponse(cachedDoc), nil
	} else if err != redis.Nil {
		fmt.Printf("Failed to fetch cached document from Redis: %v\n", err.Error())

		// Continue with execution, regardless of cache retrieval failure
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Scan DynamoDB table to retrieve ALL documents
	tableName := os.Getenv("AWS_IMAGE_TABLE")

	expr, err := util.BuildFilterConditions(filters)
	if err != nil {
		return util.InternalServerError(err), err
	}

	// Perform the scan with any conditions that may be present
	scanInput := &dynamodb.ScanInput{
		TableName:        &tableName,
		FilterExpression: expr.Condition(),
	}

	scanOutput, err := client.Scan(ctx, scanInput)
	if err != nil {
		return util.InternalServerError(err), err
	}

	// Convert response items to list of ImageDocument
	documents := []util.ImageDocument{}

	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &documents)
	if err != nil {
		return util.InternalServerError(err), err
	}

	// Sort documents by DateCreated, with latest first
	util.SortDocuments(documents)

	// Convert documents to JSON
	response, err := json.Marshal(documents)
	if err != nil {
		return util.InternalServerError(err), err
	}

	// Convert JSON to string
	responseStr := string(response)

	util.CacheJSONString(ctx, redisKey, responseStr, filters, redisClient)

	return util.JSONStringResponse(responseStr), nil
}
