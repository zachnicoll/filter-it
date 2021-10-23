package main

import (
	"aws-scalable-image-filter/internal/pkg/lambda_queue"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(lambda_queue.HandleRequest)
}
