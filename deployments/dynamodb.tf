variable "dynamodb_name" {
  default = "filterit-documents"
}

// TODO: Ensure attributes line up with Documentation in Notion
resource "aws_dynamodb_table" "ddbtable" {
  name           = var.dynamodb_name
  hash_key       = "id"
  range_key      = "date_created"
  read_capacity  = 20
  write_capacity = 20

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
    name = "filters"

    // Will be a list of filter enumerations
    type = "N"
  }

  # Make the filters array an index, so we can query by filter
  global_secondary_index {
    name            = "ImageFiltersIndex"
    hash_key        = "filters"
    range_key       = "date_created"
    read_capacity   = 20
    write_capacity  = 20
    projection_type = "ALL"
  }

  lifecycle {
    prevent_destroy = true
  }
}
