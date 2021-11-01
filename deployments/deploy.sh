#!/bin/bash
sudo apt-get remove docker docker-engine docker.io
sudo apt-get update -qq
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

sudo apt-get update
sudo apt-get install docker-ce -y

sudo usermod -a -G docker ubuntu

docker pull znicoll/filter-it-image-processor

sudo touch /etc/systemd/system/docker.filterit.service
echo "[Unit]
Description=FilterIt Image Processor Container
Requires=docker.service
After=docker.service

[Service]
TimeoutStartSec=0
Restart=always
ExecStart=/usr/bin/docker run \
  -e S3_BUCKET=${AWS_S3_BUCKET} \
  -e AWS_IMAGE_TABLE=${AWS_TABLE} \
  -e AWS_SQS_QUEUE=${AWS_SQS} \
  -e AWS_REDIS_ADDRESS=${AWS_REDIS} \
  -e AWS_REGION=${AWS_REGION} \
  znicoll/filter-it-image-processor

[Install]
WantedBy=multi-user.target" | sudo dd of=/etc/systemd/system/docker.filterit.service

sudo systemctl enable docker.filterit.service
sudo systemctl start docker.filterit.service