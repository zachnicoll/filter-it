package main

import (
	"aws-scalable-image-filter/internal/pkg/lambda_feed"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(lambda_feed.HandleRequest)
}
