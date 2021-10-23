package lambda_feed

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type FeedRequestBody struct {
	Filters []string `json:"filters"`
}

func HandleRequest(ctx context.Context, event FeedRequestBody) (*events.APIGatewayProxyResponse, error) {

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: strings.Join(event.Filters, " "),
	}, nil
}
