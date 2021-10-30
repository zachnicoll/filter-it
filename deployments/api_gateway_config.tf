variable "gateway_name" {
  default = "filterit-endpoint"
}


resource "aws_api_gateway_rest_api" "lambda" {
  name = var.gateway_name
}

resource "aws_api_gateway_deployment" "deployment" {
  depends_on = [
    aws_api_gateway_method.feedGetMethod,
    aws_api_gateway_method.queuePostMethod,
    aws_api_gateway_method.uploadPostMethod,
    aws_api_gateway_integration.queueIntegration,
    aws_api_gateway_integration.uploadIntegration,
    aws_api_gateway_integration.feedIntegration
  ]

  rest_api_id = aws_api_gateway_rest_api.lambda.id

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "prod" {
  stage_name    = "prod"
  rest_api_id   = aws_api_gateway_rest_api.lambda.id
  deployment_id = aws_api_gateway_deployment.deployment.id
}
