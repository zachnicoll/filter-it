/* Queue Lambada API Gateway Integration */
resource "aws_api_gateway_resource" "progressResource" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  parent_id   = aws_api_gateway_rest_api.lambda.root_resource_id
  path_part   = "progress"
}

resource "aws_api_gateway_method" "progressPostMethod" {
  rest_api_id   = aws_api_gateway_rest_api.lambda.id
  resource_id   = aws_api_gateway_resource.progressResource.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "progressResponse" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.progressResource.id
  http_method = aws_api_gateway_method.progressPostMethod.http_method
  status_code = "200"
}

resource "aws_api_gateway_integration" "progressIntegration" {
  rest_api_id             = aws_api_gateway_rest_api.lambda.id
  resource_id             = aws_api_gateway_resource.progressResource.id
  http_method             = aws_api_gateway_method.progressPostMethod.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambda_progress.invoke_arn
}

resource "aws_api_gateway_integration_response" "progressIntResponse" {
  depends_on = [aws_api_gateway_integration.progressIntegration]

  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.progressResource.id
  http_method = aws_api_gateway_method.progressPostMethod.http_method
  status_code = aws_api_gateway_method_response.progressResponse.status_code
}


resource "aws_api_gateway_method_settings" "progressMethod" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  stage_name  = aws_api_gateway_stage.prod.stage_name

  method_path = "*/*"

  settings {
    throttling_rate_limit  = 5
    throttling_burst_limit = 10
  }
}

resource "aws_lambda_permission" "progress_lambda_permission" {
  statement_id  = "AllowQueueAPIInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_progress.function_name
  principal     = "apigateway.amazonaws.com"

  # The /*/*/* part allows invocation from any stage, method and resource path
  # within API Gateway REST API.
  source_arn = "${aws_api_gateway_rest_api.lambda.execution_arn}/*/*/*"
}
