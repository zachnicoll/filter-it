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

  request_parameters = {
    "method.request.querystring.filter" = false
  }
}

resource "aws_api_gateway_method_response" "feedResponse" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.feedResource.id
  http_method = aws_api_gateway_method.feedGetMethod.http_method
  status_code = "200"
}

resource "aws_api_gateway_integration" "feedIntegration" {
  rest_api_id             = aws_api_gateway_rest_api.lambda.id
  resource_id             = aws_api_gateway_resource.feedResource.id
  http_method             = aws_api_gateway_method.feedGetMethod.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambda_feed.invoke_arn
}

resource "aws_api_gateway_integration_response" "feedIntResponse" {
  depends_on = [aws_api_gateway_integration.feedIntegration]

  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.feedResource.id
  http_method = aws_api_gateway_method.feedGetMethod.http_method
  status_code = aws_api_gateway_method_response.feedResponse.status_code
}

resource "aws_api_gateway_method_settings" "feedMethodSettings" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  stage_name  = aws_api_gateway_stage.prod.stage_name

  method_path = "*/*"

  settings {
    throttling_rate_limit  = 5
    throttling_burst_limit = 10
  }
}

resource "aws_lambda_permission" "feed_lambda_permission" {
  statement_id  = "AllowFeedAPIInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_feed.function_name
  principal     = "apigateway.amazonaws.com"

  # The /*/*/* part allows invocation from any stage, method and resource path
  # within API Gateway REST API.
  source_arn = "${aws_api_gateway_rest_api.lambda.execution_arn}/*/*/*"
}
