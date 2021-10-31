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

export AWS_IMAGE_TABLE=${AWS_TABLE}
export AWS_SQS_QUEUE=${AWS_SQS}
export AWS_REDIS_ADDRESS=${AWS_REDIS}
export S3_BUCKET=${AWS_S3_BUCKET}

sudo touch /etc/systemd/system/docker.filterit.service
echo "[Unit]
Description=FilterIt Image Processor Container
Requires=docker.service
After=docker.service

[Service]
TimeoutStartSec=0
Restart=always
Environment=\"S3_BUCKET=$(echo $S3_BUCKET)\"
Environment=\"AWS_IMAGE_TABLE=$(echo $AWS_IMAGE_TABLE)\"
Environment=\"AWS_SQS_QUEUE=$(echo $AWS_SQS_QUEUE)\"
Environment=\"AWS_REDIS_ADDRESS=$(echo $AWS_REDIS_ADDRESS)\"
Environment=\"AWS_ACCESS_KEY_ID=$(echo $AWS_KEY_ID)\"
Environment=\"AWS_SECRET_ACCESS_KEY=$(echo $AWS_SECRET_ACCESS_KEY)\"
Environment=\"AWS_SESSION_TOKEN=$(echo $AWS_SESSION_TOKEN)\"
ExecStartPre=-/usr/bin/docker exec %n stop
ExecStartPre=-/usr/bin/docker rm %n
ExecStartPre=/usr/bin/docker pull znicoll/filter-it-image-processor
ExecStart=/usr/bin/docker run --rm --name %n \
  -e S3_BUCKET=$(echo $S3_BUCKET) \
  -e AWS_IMAGE_TABLE=$(echo $AWS_IMAGE_TABLE) \
  -e AWS_SQS_QUEUE=$(echo $AWS_SQS_QUEUE) \
  -e AWS_REDIS_ADDRESS=$(echo $AWS_REDIS_ADDRESS) \
  -e AWS_ACCESS_KEY_ID=$(echo $AWS_KEY_ID) \
  -e AWS_SECRET_ACCESS_KEY=$(echo $AWS_SECRET_ACCESS_KEY) \
  -e AWS_SESSION_TOKEN=$(echo $AWS_SESSION_TOKEN) \
  znicoll/filter-it-image-processor

[Install]
WantedBy=multi-user.target" | sudo dd of=/etc/systemd/system/docker.filterit.service

sudo systemctl enable docker.filterit.service
sudo systemctl start docker.filterit.service