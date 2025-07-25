terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-west-2"
  default_tags {
    tags = {
      "Terraform" = "true"
      "App"       = "serve"
    }
  }
}

provider "aws" {
  alias   = "dns"
  region  = "us-east-1"
  profile = "admin"
}