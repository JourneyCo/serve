data "aws_availability_zones" "available" {}

data "aws_acm_certificate" "serve" {
  domain   = "serve.ravn.systems"
  statuses = ["ISSUED"]
}

data "aws_vpc" "this" {
  filter {
    name = "tag:Name"
    values = ["serve-uw2"]
  }
}

data "aws_subnets" "public" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.this.id]
  }

  tags = {
    Name = "serve-uw2-public-*"
  }
}
