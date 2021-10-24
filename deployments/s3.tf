variable "image_bucket" {
  default = "filterit-images"
}

variable "website_bucket" {
  default = "filterit-static-site"
}

resource "aws_s3_bucket" "image_bucket" {
  bucket = var.image_bucket

  # lifecycle {
  #   prevent_destroy = true
  # }
}

resource "aws_s3_bucket" "website_bucket" {
  bucket = var.website_bucket
  acl    = "public-read"
  policy = data.aws_iam_policy_document.website_policy.json
  website {
    index_document = "index.html"
    error_document = "404.html"
  }
}

resource "aws_s3_bucket_object" "website_objects" {
  depends_on = [
    null_resource.compile_website
  ]

  for_each = fileset("./out/website", "*")
  bucket   = aws_s3_bucket.website_bucket.id
  key      = each.value
  source   = "./out/website/${each.value}"
  etag     = filemd5("./out/website/${each.value}")
}
