data "aws_availability_zones" "available" {}

data "aws_acm_certificate" "serve" {
  domain   = "serve.ravn.systems"
  statuses = ["ISSUED"]
}
