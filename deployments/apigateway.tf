variable "gateway_name" {
  default = "filterit-endpoint"
}


resource "aws_api_gateway_rest_api" "lambda" {
  name = var.gateway_name
}

resource "aws_api_gateway_resource" "uploadResource" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  parent_id   = aws_api_gateway_rest_api.lambda.root_resource_id
  path_part   = "upload"
}

resource "aws_api_gateway_method" "uploadLambda" {
  rest_api_id   = aws_api_gateway_rest_api.lambda.id
  resource_id   = aws_api_gateway_resource.uploadResource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "uploadResponse" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.uploadResource.id
  http_method = aws_api_gateway_method.uploadLambda.http_method
  status_code = "200"
}

resource "aws_api_gateway_integration" "uploadIntegration" {
  rest_api_id             = aws_api_gateway_rest_api.lambda.id
  resource_id             = aws_api_gateway_resource.uploadResource.id
  http_method             = aws_api_gateway_method.uploadLambda.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambda_upload.invoke_arn
}

resource "aws_api_gateway_integration_response" "uploadIntResponse" {
  depends_on = [aws_api_gateway_integration.uploadIntegration]

  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.uploadResource.id
  http_method = aws_api_gateway_method.uploadLambda.http_method
  status_code = aws_api_gateway_method_response.uploadResponse.status_code
}

resource "aws_api_gateway_stage" "uploadStage" {
  stage_name    = "prod"
  rest_api_id   = aws_api_gateway_rest_api.lambda.id
  deployment_id = aws_api_gateway_deployment.uploadDeployment.id
}

resource "aws_api_gateway_deployment" "uploadDeployment" {
  depends_on = [
    aws_api_gateway_integration_response.uploadIntResponse,
    aws_api_gateway_method_response.uploadResponse,
  ]
  rest_api_id = aws_api_gateway_rest_api.lambda.id
}

resource "aws_api_gateway_method_settings" "uploadMethod" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  stage_name  = aws_api_gateway_stage.uploadStage.stage_name

  method_path = "*/*"

  settings {
    throttling_rate_limit  = 5
    throttling_burst_limit = 10
  }
}

/* Queue Lambada API Gateway Integration */
resource "aws_api_gateway_resource" "queueResource" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  parent_id   = aws_api_gateway_rest_api.lambda.root_resource_id
  path_part   = "queue"
}

resource "aws_api_gateway_method" "queueLambda" {
  rest_api_id   = aws_api_gateway_rest_api.lambda.id
  resource_id   = aws_api_gateway_resource.queueResource.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "queueResponse" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.queueResource.id
  http_method = aws_api_gateway_method.queueLambda.http_method
  status_code = "200"
}

resource "aws_api_gateway_integration" "queueIntegration" {
  rest_api_id             = aws_api_gateway_rest_api.lambda.id
  resource_id             = aws_api_gateway_resource.queueResource.id
  http_method             = aws_api_gateway_method.queueLambda.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambda_queue.invoke_arn
}

resource "aws_api_gateway_integration_response" "queueIntResponse" {
  depends_on = [aws_api_gateway_integration.queueIntegration]

  rest_api_id = aws_api_gateway_rest_api.lambda.id
  resource_id = aws_api_gateway_resource.queueResource.id
  http_method = aws_api_gateway_method.queueLambda.http_method
  status_code = aws_api_gateway_method_response.queueResponse.status_code
}

resource "aws_api_gateway_stage" "queueStage" {
  stage_name    = "prod"
  rest_api_id   = aws_api_gateway_rest_api.lambda.id
  deployment_id = aws_api_gateway_deployment.queueDeployment.id
}

resource "aws_api_gateway_deployment" "queueDeployment" {
  depends_on = [
    aws_api_gateway_integration_response.queueIntResponse,
    aws_api_gateway_method_response.queueResponse,
  ]
  rest_api_id = aws_api_gateway_rest_api.lambda.id
}

resource "aws_api_gateway_method_settings" "queueMethod" {
  rest_api_id = aws_api_gateway_rest_api.lambda.id
  stage_name  = aws_api_gateway_stage.queueStage.stage_name

  method_path = "*/*"

  settings {
    throttling_rate_limit  = 5
    throttling_burst_limit = 10
  }
}
