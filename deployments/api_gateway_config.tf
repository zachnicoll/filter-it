variable "gateway_name" {
  default = "filterit-endpoint"
}


resource "aws_api_gateway_rest_api" "lambda" {
  name = var.gateway_name
}
