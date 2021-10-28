package util

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func InternalServerError(err error, method string) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods": method,
		},
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}
}

func JSONStringResponse(body string, method string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods": method,
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
func BuildFilterConditions(filters []int) (expression.Expression, error) {
	conditions := []expression.ConditionBuilder{}

	for _, filter := range filters {
		condition := expression.Name("filters").In(expression.Value(filter))
		conditions = append(conditions, condition)
	}

	builder := expression.NewBuilder()

	for _, condition := range conditions {
		builder = builder.WithCondition(condition)
	}

	// Make sure that only DONE documents are selected from DynamoDB
	progressCondition := expression.Name("progress").Equal(expression.Value(DONE))
	builder = builder.WithCondition(progressCondition)

	return builder.Build()
}
