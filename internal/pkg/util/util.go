package util

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func InternalServerError(err error) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods":     "GET, PUT, PATCH, POST, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Authorization, Content-Type",
		},
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}
}

func JSONStringResponse(body string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods":     "GET, PUT, PATCH, POST, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Authorization, Content-Type",
		},
		Body: body,
	}
}

func SortFilters(filters []int) {
	sort.Slice(
		filters,
		func(i, j int) bool {
			return i > j
		},
	)
}

/*
Sort ImageDocument slice by DateCreated, with latest first.
*/
func SortDocuments(documents []ImageDocument) {
	sort.Slice(
		documents,
		func(i, j int) bool {
			return documents[i].DateCreated < documents[j].DateCreated
		},
	)
}

/*
Create a DynamoDB Expression containing conditions where each document retrieved
must include each of the supplied filters in the "filters" attribute.
*/
func BuildFilterConditions(filter string) (expression.Expression, error) {
	builder := expression.NewBuilder()

	filterInt, err := strconv.Atoi(filter)

	if err == nil {
		filterCondition := expression.Name("filter").Equal(expression.Value(filterInt))
		builder = builder.WithCondition(filterCondition)
	}

	// Make sure that only DONE documents are selected from DynamoDB
	progressCondition := expression.Name("progress").Equal(expression.Value(DONE))

	builder = builder.WithCondition(progressCondition)

	return builder.Build()
}

func FetchInstanceID() string {
	response, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		log.Fatalln("Unable to find instance ID")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("Unable to close instance ID response")
		}
	}(response.Body)

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("Unable to read instance ID response")
	}

	instanceID := string(responseData)

	return instanceID
}

/*
	Updates the given document in DynamoDB, and invalidates the Redis entry for its
	corresponding Tag (and the "ALL" key, "_").
*/
func UpdateDocument(ctx context.Context, clients *Clients, table *string, item *ImageDocument) (err error) {
	itemMap, err := attributevalue.MarshalMap(item)

	if err != nil {
		return err
	}
	_, err = clients.DynamoDb.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      itemMap,
		TableName: table,
	})

	if err != nil {
		return err
	}

	log.Println("Invalidating cache key", item.Tag)

	InvalidateCache(ctx, fmt.Sprintf("%d", item.Tag), clients.Redis)

	return nil
}
