terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

// TODO: configure with uni config
provider "aws" {
  access_key                  = "mock_access_key"
  region                      = "ap-southeast-2"
  secret_key                  = "mock_secret_key"
  s3_force_path_style         = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway = "http://localhost:4566"
    cloudwatch = "http://localhost:4566"
    dynamodb   = "http://localhost:4566"
    iam        = "http://localhost:4566"
    lambda     = "http://localhost:4566"
    s3         = "http://localhost:4566"
    sqs        = "http://localhost:4566"
  }
}
