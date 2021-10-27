package lambda_feed

import (
	"aws-scalable-image-filter/internal/pkg/helpers"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-redis/redis/v8"
)

func HandleRequest(ctx context.Context, event helpers.FeedRequestBody) *events.APIGatewayProxyResponse {
	// Load default AWS config, including AWS_REGION env var
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return helpers.InternalServerError(err)
	}

	// Sort filters
	filters := event.Filters
	sort.Slice(
		filters,
		func(i, j int) bool {
			return i < j
		},
	)

	// Connect to Redis
	redisHost := os.Getenv("AWS_REDIS_ADDRESS")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", redisHost),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create a key that is the concatenation of all filters in ascending order.
	// E.g. [2, 1, 3] creates the key "123"
	redisKey := "_" // "_" is the "root" key, where a search for ALL documents is stored
	for _, filter := range filters {
		redisKey = fmt.Sprintf("%s%d", redisKey, filter)
	}

	cachedDoc, err := rdb.Get(ctx, redisKey).Result()

	// Found a cached document for this query, return it
	if err == nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: cachedDoc,
		}
	} else if err != nil {
		fmt.Printf("Failed to fetch cached document from Redis: %v", err.Error())

		// Continue with execution, regardless of cache retrieval failure
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Scan DynamoDB table to retrieve ALL documents
	tableName := os.Getenv("AWS_IMAGE_TABLE")

	// Create a condition for each filter supplied with the request,
	// so that each retrieved document must include each of the supplied filters
	conditions := []expression.ConditionBuilder{}

	for _, filter := range filters {
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

	resp, err := client.Scan(ctx, input)
	if err != nil {
		return helpers.InternalServerError(err)
	}

	// Convert response items to list of ImageDocument
	documents := []helpers.ImageDocument{}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &documents)
	if err != nil {
		return helpers.InternalServerError(err)
	}

	// Sort documents by DateCreated, with latest first
	sort.Slice(
		documents,
		func(i, j int) bool {
			return documents[i].DateCreated > documents[j].DateCreated
		},
	)

	// Convert documents to JSON
	response, err := json.Marshal(documents)
	if err != nil {
		return helpers.InternalServerError(err)
	}

	// Convert JSON to string
	responseStr := string(response)

	// Store JSON string in Redis
	_, err = rdb.Set(ctx, redisKey, responseStr, 0).Result()
	if err != nil {
		fmt.Printf("Failed to fetch cache document in Redis: %v", err.Error())
	}

	/**
	Add Redis key (e.g. "_123") to each index Set.
	See explanation of index Sets in Notion:
	https://www.notion.so/Caching-Strategy-8b25c9a2f1354186b99e1e313e45d3ae
	*/
	for _, filter := range filters {
		// Format index like "x:index", where x is a filter e.g. "3:index"
		_, err = rdb.SAdd(ctx, fmt.Sprintf("%d:index", filter), redisKey).Result()
		if err != nil {
			fmt.Printf("Failed to add key to Set index in Redis: %v", err.Error())
		}
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}
}
