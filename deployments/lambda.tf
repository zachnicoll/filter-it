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
    }
  }
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
      AWS_IMAGE_TABLE   = var.dynamodb_name
      S3_BUCKET = var.image_bucket
    }
  }

  vpc_config {
    subnet_ids         = [aws_subnet.filterit-subnet-private-1.id]
    security_group_ids = [aws_security_group.filterit-sg.id]
  }
}
