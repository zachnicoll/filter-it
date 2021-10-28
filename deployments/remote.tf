resource "null_resource" "remote_image_processor" {
  depends_on = [
    null_resource.dockerise_image_processor
  ]

  connection {
    type        = "ssh"
    user        = "ec2-user"
    private_key = file(var.PEM_PATH)
    host        = aws_instance.image_processor.public_ip
  }

  provisioner "remote-exec" {
    inline = [
      "sudo yum update -y",
      "sudo yum install docker -y",
      "sudo service docker start",
      "sudo usermod -a -G docker ec2-user",
      "docker pull znicoll/filter-it-image-processor",
      "docker run znicoll/filter-it-image-processor"
    ]
  }
}
