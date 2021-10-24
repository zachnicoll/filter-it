variable "image_bucket" {
  default = "filterit-images"
}

variable "website_bucket" {
  default = "filterit-staticsite"
}

resource "aws_s3_bucket" "image_bucket" {
  bucket = var.image_bucket

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_s3_bucket" "website_bucket" {
  bucket = var.website_bucket
  acl    = "public-read"
  policy = data.aws_iam_policy_document.website_policy.json

  website {
    index_document = "index.html"
    error_document = "404.html"
  }

  force_destroy = true
  lifecycle {
    prevent_destroy = false
  }
}

# TODO: Get auto uploading of website working
# resource "aws_s3_bucket_object" "website_objects" {
#   depends_on = [
#     null_resource.compile_website
#   ]

#   for_each = fileset("${path.module}/out/website/", "**")

#   bucket = aws_s3_bucket.website_bucket.id
#   key    = each.value
#   source = "${path.module}/out/website/${each.value}"
#   etag   = filemd5("./out/website/${each.value}")
# }
