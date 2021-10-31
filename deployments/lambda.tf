/* Lambada Upload */
resource "aws_lambda_function" "lambda_upload" {
  filename         = data.archive_file.lambda_upload_zip.output_path
  handler          = "lambda_upload"
  role             = aws_iam_role.lambda_exec_role.arn
  runtime          = "go1.x"
  function_name    = "filterit-lambda_upload"
  source_code_hash = data.archive_file.lambda_upload_zip.output_base64sha256

  vpc_config {
    subnet_ids         = [aws_subnet.filterit-subnet-private-1.id]
    security_group_ids = [aws_security_group.filterit-sg.id]
  }

  environment {
    variables = {
      S3_BUCKET = var.image_bucket
    }
  }
}

output "upload_output" {
  value       = "${aws_api_gateway_stage.prod.invoke_url}${aws_api_gateway_resource.uploadResource.path}"
  description = "The public upload link for invoking upload gateway"
}

/* Lambda Queue */
resource "aws_lambda_function" "lambda_queue" {
  filename         = data.archive_file.lambda_queue_zip.output_path
  handler          = "lambda_queue"
  role             = aws_iam_role.lambda_exec_role.arn
  runtime          = "go1.x"
  function_name    = "filterit-lambda_queue"
  source_code_hash = data.archive_file.lambda_queue_zip.output_base64sha256

  environment {
    variables = {
      AWS_IMAGE_TABLE   = var.dynamodb_name
      AWS_SQS_QUEUE     = var.sqs_name
      AWS_REDIS_ADDRESS = element(aws_elasticache_cluster.redis.cache_nodes, 0).address
    }
  }

  vpc_config {
    subnet_ids         = [aws_subnet.filterit-subnet-private-1.id]
    security_group_ids = [aws_security_group.filterit-sg.id]
  }
}

output "queue_output" {
  value       = "${aws_api_gateway_stage.prod.invoke_url}${aws_api_gateway_resource.queueResource.path}"
  description = "The public queue link for invoking upload gateway"
}

/* Lambda Feed*/
resource "aws_lambda_function" "lambda_feed" {
  filename         = data.archive_file.lambda_feed_zip.output_path
  handler          = "lambda_feed"
  role             = aws_iam_role.lambda_exec_role.arn
  runtime          = "go1.x"
  function_name    = "filterit-lambda_feed"
  source_code_hash = data.archive_file.lambda_feed_zip.output_base64sha256

  vpc_config {
    subnet_ids         = [aws_subnet.filterit-subnet-private-1.id]
    security_group_ids = [aws_security_group.filterit-sg.id]
  }

  environment {
    variables = {
      AWS_IMAGE_TABLE   = var.dynamodb_name
      AWS_REDIS_ADDRESS = element(aws_elasticache_cluster.redis.cache_nodes, 0).address
      S3_BUCKET         = var.image_bucket
    }
  }
}

output "feed_output" {
  value       = "${aws_api_gateway_stage.prod.invoke_url}${aws_api_gateway_resource.feedResource.path}"
  description = "The public feed link for invoking upload gateway"
}

/* Lambda Progress */
resource "aws_lambda_function" "lambda_progress" {
  filename         = data.archive_file.lambda_progress_zip.output_path
  handler          = "lambda_progress"
  role             = aws_iam_role.lambda_exec_role.arn
  runtime          = "go1.x"
  function_name    = "filterit-lambda_progress"
  source_code_hash = data.archive_file.lambda_progress_zip.output_base64sha256

  environment {
    variables = {
      AWS_IMAGE_TABLE = var.dynamodb_name
      S3_BUCKET       = var.image_bucket
    }
  }

  vpc_config {
    subnet_ids         = [aws_subnet.filterit-subnet-private-1.id]
    security_group_ids = [aws_security_group.filterit-sg.id]
  }
}

output "progress_output" {
  value       = "${aws_api_gateway_stage.prod.invoke_url}${aws_api_gateway_resource.progressResource.path}"
  description = "The public progress link for invoking upload gateway"
}
