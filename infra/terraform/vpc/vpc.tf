locals {
  azs = slice(data.aws_availability_zones.available.names, 0, 3)
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.19.0"

  name = "serve-uw2"
  cidr = "10.100.0.0/16"

  azs             = ["us-west-2a", "us-west-2b", "us-west-2c"]
  private_subnets = ["10.100.1.0/24", "10.100.2.0/24", "10.100.3.0/24"]
  public_subnets  = ["10.100.101.0/24", "10.100.102.0/24", "10.100.103.0/24"]
  database_subnets = ["10.100.50.0/24", "10.100.51.0/24", "10.100.52.0/24"]

  enable_nat_gateway = false
  enable_vpn_gateway = false

  tags = {
    "Terraform" = "true"
    "App"       = "serve"
  }
  vpc_tags = {
    "Name" : "serve-uw2"
  }

}
