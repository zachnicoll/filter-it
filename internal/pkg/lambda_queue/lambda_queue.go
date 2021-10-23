package lambda_queue

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func HandleRequest(ctx context.Context) (*events.APIGatewayProxyResponse, error) {

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "",
	}, nil
}
