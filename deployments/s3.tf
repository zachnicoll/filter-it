variable "image_bucket" {
  default = "filterit-images"
}

resource "aws_s3_bucket" "bucket" {
  bucket = var.image_bucket

  lifecycle {
    prevent_destroy = true
  }
}
