data "aws_vpc" "main" {
  filter {
    name   = "tag:Name"
    values = ["OneTwentySeven"]
  }
}

data "aws_subnets" "dmz" {
  filter {
    name   = "tag:dmz"
    values = ["true"]
  }

  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.main.id]
  }
}


data "aws_subnets" "app" {
  filter {
    name   = "tag:app"
    values = ["true"]
  }

  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.main.id]
  }
}
