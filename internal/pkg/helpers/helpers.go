package helpers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func InternalServerError(err error) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}
}
