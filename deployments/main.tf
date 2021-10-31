terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  backend "s3" {
    bucket = "filterit-terraform-state"
    key    = "terraform.tfstate"
    region = "ap-southeast-2"
  }
}

provider "aws" {
  region = "ap-southeast-2"
}

output "invoke_output" {
  value       = aws_api_gateway_stage.prod.invoke_url
  description = "The public IP address of gateway stage"
}
