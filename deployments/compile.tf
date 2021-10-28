/* Compile Lambda Functions */

resource "null_resource" "compile_lambda_feed" {
  triggers = {
    build_number = timestamp()
  }

  provisioner "local-exec" {
    command = "GOOS=linux go build -o ./out/lambda_feed ../cmd/lambda_feed"
  }
}

resource "null_resource" "compile_lambda_upload" {
  triggers = {
    build_number = timestamp()
  }

  provisioner "local-exec" {
    command = "GOOS=linux go build -o ./out/lambda_upload ../cmd/lambda_upload"
  }
}

resource "null_resource" "compile_lambda_queue" {
  triggers = {
    build_number = timestamp()
  }

  provisioner "local-exec" {
    command = "GOOS=linux go build -o ./out/lambda_queue ../cmd/lambda_queue"
  }
}

# resource "null_resource" "dockerise_image_processor" {
#   triggers = {
#     build_number = timestamp()
#   }

#   provisioner "local-exec" {
#     command = <<-EOT
#     docker build -f ../build/Dockerfile.filterit -t znicoll/filter-it-image-processor ../
#     docker push znicoll/filter-it-image-processor:latest
#     EOT
#   }
# }

# TODO: Get website build working
# resource "null_resource" "compile_website" {
#   triggers = {
#     build_number = "${timestamp()}"
#   }

#   provisioner "local-exec" {
#     command = "rm -rf ${path.module}/out/website/ && mkdir ${path.module}/out/website/ && cd ../website && yarn && yarn build && cd - && mv ../website/out/* ${path.module}/out/website/"
#   }
# }

/* Archiver Executables */
data "archive_file" "lambda_feed_zip" {
  source_file = "${path.module}/out/lambda_feed"
  output_path = "${path.module}/out/lambda_feed.zip"
  type        = "zip"
  depends_on = [
    null_resource.compile_lambda_feed
  ]
}

data "archive_file" "lambda_queue_zip" {
  source_file = "${path.module}/out/lambda_queue"
  output_path = "${path.module}/out/lambda_queue.zip"
  type        = "zip"
  depends_on = [
    null_resource.compile_lambda_queue
  ]
}

data "archive_file" "lambda_upload_zip" {
  source_file = "${path.module}/out/lambda_upload"
  output_path = "${path.module}/out/lambda_upload.zip"
  type        = "zip"
  depends_on = [
    null_resource.compile_lambda_upload
  ]
}
