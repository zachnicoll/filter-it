package lambda_feed

import (
	"aws-scalable-image-filter/internal/pkg/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-redis/redis/v8"
)

func HandleRequest(_ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx := context.Background()

	// Load default AWS config, including AWS_REGION env var
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return util.InternalServerError(err), err
	}
	filter := request.QueryStringParameters["filter"]

	// If filter is not present in query, set it to the "all" redis key
	if filter == "" {
		filter = "_"
	}

	redisClient := util.ConnectToRedis()

	redisKey := filter
	cachedDoc, err := redisClient.Get(ctx, redisKey).Result()

	redisCacheHit := false

	if err == nil && cachedDoc != "" {
		// Found a cached document for this query, return it
		redisCacheHit = true
		// return util.JSONStringResponse(cachedDoc), nil
	} else if err != redis.Nil {
		fmt.Printf("Failed to fetch cached document from Redis: %v\n", err.Error())

		// Continue with execution, regardless of cache retrieval failure
	}

	documents := []util.ImageDocument{}

	if !redisCacheHit {
		// Create DynamoDB client
		client := dynamodb.NewFromConfig(cfg)
		tableName := os.Getenv("AWS_IMAGE_TABLE")

		if filter == "_" {
			// Scan DynamoDB table to retrieve ALL documents

			filt := expression.Name("progress").Equal(expression.Value(util.DONE))
			expr, err := expression.NewBuilder().WithFilter(filt).Build()

			if err != nil {
				return util.InternalServerError(err), err
			}

			// Perform the scan with any conditions that may be present
			scanInput := &dynamodb.ScanInput{
				TableName:                 &tableName,
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				FilterExpression:          expr.Filter(),
				ProjectionExpression:      expr.Projection(),
			}

			scanOutput, err := client.Scan(ctx, scanInput)
			if err != nil {
				return util.InternalServerError(err), err
			}

			// Convert response items to list of ImageDocument
			err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &documents)
			if err != nil {
				return util.InternalServerError(err), err
			}
		} else {
			indexName := "tag-date_created-index"
			queryInput := &dynamodb.QueryInput{
				TableName: &tableName,
				IndexName: &indexName,
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":f": &types.AttributeValueMemberN{Value: filter},
				},
				KeyConditionExpression: aws.String("tag = :f"),
			}

			queryOutput, err := client.Query(ctx, queryInput)
			if err != nil {
				return util.InternalServerError(err), err
			}

			// Convert response items to list of ImageDocument
			err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &documents)
			if err != nil {
				return util.InternalServerError(err), err
			}
		}
	} else {
		// Convert cached Redis JSON string to ImageDocument[]
		err := json.Unmarshal([]byte(cachedDoc), &documents)
		if err != nil {
			return util.InternalServerError(err), err
		}
	}

	// Sort documents by DateCreated, with latest first
	util.SortDocuments(documents)

	s3BucketName := os.Getenv("S3_BUCKET")
	if s3BucketName == "" {
		return util.InternalServerError(errors.New("S3 Bucket was unable to be loaded from env vars")), nil
	}

	s3Client := s3.NewFromConfig(cfg)
	s3PresignClient := s3.NewPresignClient(s3Client)

	// Always attach a new signed URL, even for cached results
	signedDocuments := []util.ImageDocument{}

	for _, doc := range documents {
		input := &s3.GetObjectInput{
			Bucket: &s3BucketName,
			Key:    &doc.Image,
		}

		resp, err := s3PresignClient.PresignGetObject(ctx, input)
		if err != nil {
			return util.InternalServerError(err), err
		}

		signedUrl := url.QueryEscape(resp.URL)
		doc.ImageURL = signedUrl

		signedDocuments = append(signedDocuments, doc)
	}

	// Convert documents to JSON
	response, err := json.Marshal(signedDocuments)
	if err != nil {
		return util.InternalServerError(err), err
	}

	// Convert JSON to string
	responseStr := string(response)

	util.CacheJSONString(ctx, redisKey, responseStr, redisClient)

	return util.JSONStringResponse(responseStr), nil
}
