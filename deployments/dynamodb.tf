variable "dynamodb_name" {
  default = "filterit-documents"
}

// TODO: Ensure attributes line up with Documentation in Notion
resource "aws_dynamodb_table" "ddbtable" {
  name           = var.dynamodb_name
  hash_key       = "id"
  range_key      = "date_created"
  read_capacity  = 5
  write_capacity = 5

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "date_created"

    // UNIX timestamp (seconds)
    type = "N"
  }

  attribute {
    name = "filter"

    // Will be a list of filter enumerations
    type = "N"
  }

  # Make the filters array an index, so we can query by filter
  global_secondary_index {
    name            = "filter-date_create-index"
    hash_key        = "filter"
    range_key       = "date_created"
    read_capacity   = 5
    write_capacity  = 5
    projection_type = "ALL"
  }

  lifecycle {
    prevent_destroy = true
  }
}
