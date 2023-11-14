resource "aws_security_group" "alb" {
  name        = "upload_alb"
  description = "Allow inbound traffic to alb"
}

resource "aws_security_group_rule" "alb_http" {
  security_group_id = aws_security_group.alb.id
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "alb_https" {
  security_group_id = aws_security_group.alb.id
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "alb_egress" {
  security_group_id = aws_security_group.alb.id
  type              = "ingress"
  from_port         = 0
  to_port           = 65535
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "random_string" "tg_name" {
  length  = 4
  special = false

}

resource "aws_lb_target_group" "upload" {
  name     = "upload-tg-${random_string.tg_name.result}"
  protocol = "HTTP"
  port     = 3000
  vpc_id   = data.aws_vpc.main.id

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    interval            = 30
    path                = "/health"
    port                = 3000
    protocol            = "HTTP"
  }
}

resource "aws_lb" "upload" {
  name               = "upload-service-lb"
  load_balancer_type = "application"
  security_groups    = [aws_security_group.upload_service.id]
  subnets            = data.aws_subnets.dmz.ids
}



resource "aws_lb_listener" "upload_http" {
  load_balancer_arn = aws_lb.upload.arn
  port              = "80"
  protocol          = "HTTP"
  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb_listener" "upload_https" {
  load_balancer_arn = aws_lb.upload.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate.photos.arn
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.upload.arn
  }
}
