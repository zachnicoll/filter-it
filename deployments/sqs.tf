variable "sqs_name" {
  default = "filterit-queue"
}

resource "aws_sqs_queue" "queue" {
  name = var.sqs_name
}
