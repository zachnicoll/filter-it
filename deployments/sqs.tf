variable "sqs_name" {
  default = "filterit-queue.fifo"
}

resource "aws_sqs_queue" "filterit-queue" {
  name                        = var.sqs_name
  fifo_queue                  = true
  content_based_deduplication = true
}
