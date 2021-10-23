/* Feed Lambda API Gateway Integration */
resource "aws_api_gateway_resource" "feedResource" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  parent_id   = aws_api_gateway_rest_api.lambda.root_resource_id
  path_part   = "feed"
}

// GET /feed
resource "aws_api_gateway_method" "feedGetMethod" {
  rest_api_id   = aws_api_gateway_rest_api.lambda.id
  resource_id   = aws_api_gateway_resource.feedResource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "feedGetResponse" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.feedResource.id
  http_method = aws_api_gateway_method.feedGetMethod.http_method
  status_code = "200"
}

resource "aws_api_gateway_integration" "feedGetIntegration" {
  rest_api_id             = aws_api_gateway_rest_api.lambda.id
  resource_id             = aws_api_gateway_resource.feedResource.id
  http_method             = aws_api_gateway_method.feedGetMethod.id
  integration_http_method = aws_api_gateway_method.feedGetMethod.http_method
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambda_queue.invoke_arn
}

resource "aws_api_gateway_integration_response" "feedGetIntResponse" {
  depends_on = [aws_api_gateway_integration.feedGetIntegration]

  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.feedResource.id
  http_method = aws_api_gateway_method.feedGetMethod.http_method
  status_code = aws_api_gateway_method_response.feedGetResponse.status_code
}

resource "aws_api_gateway_stage" "feedStage" {
  stage_name    = "prod"
  rest_api_id   = aws_api_gateway_rest_api.lambda.id
  deployment_id = aws_api_gateway_deployment.feedDeployment.id
}

resource "aws_api_gateway_deployment" "feedDeployment" {
  depends_on = [
    aws_api_gateway_integration_response.feedGetIntResponse,
    aws_api_gateway_method_response.feedGetResponse,
  ]
  rest_api_id = aws_api_gateway_rest_api.lambda.id
}

resource "aws_api_gateway_method_settings" "feedMethodSettings" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  stage_name  = aws_api_gateway_stage.feedStage.stage_name

  method_path = "*/*"

  settings {
    throttling_rate_limit  = 5
    throttling_burst_limit = 10
  }
}
