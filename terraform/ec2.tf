resource "aws_security_group" "upload_service" {
  name   = "upload-service-sg"
  vpc_id = data.aws_vpc.main.id
}

resource "aws_security_group_rule" "allow_ssh_home" {
  security_group_id = aws_security_group.upload_service.id
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = ["208.104.177.231/32"]
}

resource "aws_security_group_rule" "allow_http" {
  security_group_id = aws_security_group.upload_service.id
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "allow_https" {
  security_group_id = aws_security_group.upload_service.id
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "allow_egress" {
  security_group_id = aws_security_group.upload_service.id
  type              = "egress"
  from_port         = 0
  to_port           = 65535
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "vpc_ingress_ssh" {
  security_group_id = aws_security_group.upload_service.id
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = [data.aws_vpc.main.cidr_block]
}

resource "aws_security_group_rule" "vpc_ingress_http" {
  security_group_id = aws_security_group.upload_service.id
  type              = "ingress"
  from_port         = 8080
  to_port           = 8080
  protocol          = "tcp"
  cidr_blocks       = [data.aws_vpc.main.cidr_block]
}


resource "aws_network_interface" "upload" {
  subnet_id       = data.aws_subnets.app.ids[0]
  security_groups = [aws_security_group.upload_service.id]
}

data "aws_key_pair" "ddouglas" {
  key_name = "devoverlord@DESKTOP-BMOGNVQ"
}

resource "aws_instance" "upload" {
  ami           = "ami-016485166ec7fa705"
  instance_type = "t4g.small"

  key_name = data.aws_key_pair.ddouglas.key_name

  network_interface {
    device_index         = 0
    network_interface_id = aws_network_interface.upload.id
  }

  root_block_device {
    delete_on_termination = true
    volume_size           = 20
    volume_type           = "gp2"
  }

  iam_instance_profile = aws_iam_instance_profile.upload.name


  user_data                   = file("./templates/upload_ec2_userdata.tpl")
  user_data_replace_on_change = true

}

output "upload_service_private_ip_address" {
  value = aws_instance.upload.private_ip
}

resource "aws_iam_instance_profile" "upload" {
  name = "upload-service-iam-profile"
  role = aws_iam_role.upload_admin.name
}

resource "aws_iam_role" "upload_admin" {
  name               = "upload-web-service-admin"
  assume_role_policy = data.aws_iam_policy_document.upload_assume_role.json

  managed_policy_arns = ["arn:aws:iam::aws:policy/AdministratorAccess"]

}

data "aws_iam_policy_document" "upload_assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

output "upload_service_public_ip" {
  value = aws_instance.upload.public_ip
}
