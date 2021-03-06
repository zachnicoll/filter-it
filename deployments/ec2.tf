variable "processor_ami_name" {
  default = "image_processor_ami"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

data "template_file" "deploy_script" {
  template = file("${path.module}/deploy.sh")

  vars = {
    AWS_TABLE     = var.dynamodb_name
    AWS_SQS       = var.sqs_name
    AWS_REDIS     = element(aws_elasticache_cluster.redis.cache_nodes, 0).address
    AWS_S3_BUCKET = var.image_bucket
    AWS_REGION    = "ap-southeast-2"
  }
}

resource "aws_instance" "image_processor" {
  ami             = data.aws_ami.ubuntu.id
  instance_type   = "t2.micro"
  security_groups = [aws_security_group.filterit-sg.id]
  subnet_id       = aws_subnet.filterit-subnet-public-1.id
  key_name        = aws_key_pair.image_processor_key.key_name
  user_data       = data.template_file.deploy_script.rendered
  iam_instance_profile = aws_iam_instance_profile.ec2_iam_profile.name

  depends_on = [
    aws_dynamodb_table.ddbtable,
    aws_sqs_queue.filterit-queue,
    aws_elasticache_cluster.redis,
  ]
}

resource "aws_ami_from_instance" "image_processor_ami" {
  name               = var.processor_ami_name
  source_instance_id = aws_instance.image_processor.id
}

