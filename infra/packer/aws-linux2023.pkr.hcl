packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "al2023" {
  ami_name      = "serve-al2023-ami-{{timestamp}}"
  instance_type = "t3.medium"
  region        = "us-west-2"
  source_ami    = "ami-05572e392e80aee89"
  ssh_username  = "ec2-user"
}

build {
  name = "serve-golden-image"
  sources = [
    "source.amazon-ebs.al2023"
  ]

  provisioner "ansible" {
    playbook_file       = "playbook/serve.yml"
    keep_inventory_file = true
    use_proxy           = false
  }
}
