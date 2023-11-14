#!/bin/bash

sudo apt update
sudo apt upgrade -y
apt install zip unzip

curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

aws s3 cp s3://ddouglas-desktop/upload ./upload