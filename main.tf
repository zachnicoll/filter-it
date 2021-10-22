terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
    }
  }
}

// TODO: configure with uni config
provider "aws" {
  access_key                  = "mock_access_key"
  region                      = "us-east-1"
  secret_key                  = "mock_secret_key"
}

variable "image_bucket" {
  default = "filterit-images"
}

variable "sqs_name" {
  default = "filterit-queue"
}

variable "dynamodb_name" {
  default = "filterit-documents"
}

variable "gateway_name" {
  default = "filterit-upload-endpoint"
}


/* Resources */
resource "aws_s3_bucket" "bucket" {
  bucket = var.image_bucket

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_sqs_queue" "queue" {
  name = var.sqs_name
}

// TODO: Ensure attributes line up with Documentation in Notion
resource "aws_dynamodb_table" "ddbtable" {
  name           = var.dynamodb_name
  hash_key       = "id"
  read_capacity  = 20
  write_capacity = 20
  attribute {
    name = "id"
    type = "S"
  }
}

/* Lambda IAM Role Setup */

resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda_exec_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "lambda_exec_logging" {
  name = "lambda_exec_logging"
  path = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "lambda_exec_sqs" {
  name = "lambda_exec_sqs"
  path = "/"
  description = "IAM policy for sqs from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "sqs:SendMessage"
      ],
      "Resource": "${aws_sqs_queue.queue.arn}",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "lambda_exec_s3" {
  name = "lambda_exec_s3"
  path = "/"
  description = "IAM policy for s3 from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:PutObject"
      ],
      "Resource": "${aws_s3_bucket.bucket.arn}",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "lambda_exec_dynamodb" {
  name = "lambda_exec_dynamodb"
  path = "/"
  description = "IAM policy for dynamodb from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:DeleteItem"
      ],
      "Resource": "${aws_dynamodb_table.ddbtable.arn}",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.lambda_exec_logging.arn
}

resource "aws_iam_role_policy_attachment" "lambda_sqs" {
  role = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.lambda_exec_sqs.arn
}

resource "aws_iam_role_policy_attachment" "lambda_s3" {
  role = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.lambda_exec_s3.arn
}

resource "aws_iam_role_policy_attachment" "lambda_dynamodb" {
  role = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.lambda_exec_dynamodb.arn
}

/* Lambada Upload */
data "archive_file" "lambda_upload_zip" {
  source_file  = "${path.module}/out/lambda_upload"
  output_path = "${path.module}/out/lambda_upload.zip"
  type        = "zip"
}

resource "aws_lambda_function" "lambda_upload" {
  filename      = data.archive_file.lambda_upload_zip.output_path
  handler       = "lambda_upload"
  role          = aws_iam_role.lambda_exec_role.arn
  runtime       = "go1.x"
  function_name = "filterit-lambda_upload"

  environment {
    variables = {
      S3_BUCKET = var.image_bucket
    }
  }
}

/* Lambda Queue */
data "archive_file" "lambda_queue_zip" {
  source_file  = "${path.module}/out/lambda_queue"
  output_path = "${path.module}/out/lambda_queue.zip"
  type        = "zip"
}

resource "aws_lambda_function" "lambda_queue" {
  filename      = data.archive_file.lambda_queue_zip.output_path
  handler       = "lambda_queue"
  role          = aws_iam_role.lambda_exec_role.arn
  runtime       = "go1.x"
  function_name = "filterit-lambda_queue"

  environment {
    variables = {
      AWS_IMAGE_TABLE = var.dynamodb_name
      AWS_SQS_QUEUE = var.sqs_name
    }
  }
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
  stage_name = "prod"
  rest_api_id = aws_api_gateway_rest_api.lambda.id
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
  stage_name = aws_api_gateway_stage.uploadStage.stage_name

  method_path = "*/*"

  settings {
    throttling_rate_limit = 5
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
  stage_name = "prod"
  rest_api_id = aws_api_gateway_rest_api.lambda.id
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
  stage_name = aws_api_gateway_stage.queueStage.stage_name

  method_path = "*/*"

  settings {
    throttling_rate_limit = 5
    throttling_burst_limit = 10
  }
}