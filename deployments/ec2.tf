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

# resource "aws_instance" "image_processor" {
#   ami             = data.aws_ami.ubuntu.id
#   instance_type   = "t2.micro"
#   security_groups = [aws_security_group.filterit-sg.id]
#   subnet_id       = aws_subnet.filterit-subnet-private-1.id
#   key_name        = aws_key_pair.image_processor_key.key_name
# }

# resource "aws_ami_from_instance" "image_processor_ami" {
#   name               = var.processor_ami_name
#   source_instance_id = aws_instance.image_processor.id

#   depends_on = [
#     null_resource.remote_image_processor
#   ]
# }

