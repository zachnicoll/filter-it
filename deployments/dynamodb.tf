variable "dynamodb_name" {
  default = "filterit-documents"
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
