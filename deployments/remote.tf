output "PEM_PATH" {
  value = var.PEM_PATH
}

resource "null_resource" "remote_image_processor" {
  depends_on = [
    null_resource.dockerise_image_processor
  ]

  connection {
    type        = "ssh"
    user        = "ubuntu"
    private_key = file(var.PEM_PATH)
    host        = aws_instance.image_processor.public_ip
  }

  provisioner "remote-exec" {
    inline = [
      "sudo yum update -y",
      "sudo yum install docker -y",
      "sudo systemctl enable docker",
      "sudo systemctl start docker",
      "sudo usermod -a -G docker ec2-user",

      # Create a systemd service that will automatically start the docker container on boot
      "sudo touch /etc/systemd/system/docker.filterit.service",
      <<-EOT
      echo "[Unit]
Description=FilterIt Image Processor Container
Requires=docker.service
After=docker.service

[Service]
Restart=always
ExecStartPre=/usr/bin/docker pull znicoll/filter-it-image-processor
ExecStart=/usr/bin/docker start -a znicoll/filter-it-image-processor
ExecStop=/usr/bin/docker stop -t 2 znicoll/filter-it-image-processor

[Install]
WantedBy=multi-user.target" | sudo dd of=/etc/systemd/system/docker.filterit.service
      EOT
      ,
      "sudo systemctl enable docker.filterit.service",
      "sudo systemctl start docker.filterit.service",
      "docker pull znicoll/filter-it-image-processor:latest",
      "docker run znicoll/filter-it-image-processor:latest"
    ]
  }
}
