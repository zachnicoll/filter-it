/* Lambada Upload */
data "archive_file" "lambda_upload_zip" {
  source_file = "${path.module}/out/lambda_upload"
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
  source_file = "${path.module}/out/lambda_queue"
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
      AWS_SQS_QUEUE   = var.sqs_name
    }
  }
}
