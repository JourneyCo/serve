resource "aws_security_group" "serve_vm" {
  name        = "serve"
  description = "server app"
  vpc_id      = data.aws_vpc.this.id
}

resource "aws_vpc_security_group_ingress_rule" "allow_lb" {
  security_group_id = aws_security_group.serve_vm.id
  cidr_ipv4         = data.aws_vpc.this.cidr_block
  ip_protocol       = "-1"
}

resource "aws_vpc_security_group_ingress_rule" "ssh" {
  security_group_id = aws_security_group.serve_vm.id
  cidr_ipv4         = "76.131.61.220/32"
  from_port         = 22
  ip_protocol       = "tcp"
  to_port           = 22
}

resource "aws_vpc_security_group_ingress_rule" "http" {
  security_group_id = aws_security_group.serve_vm.id
  cidr_ipv4         = "76.131.61.220/32"
  from_port         = 80
  ip_protocol       = "tcp"
  to_port           = 80
}

resource "aws_vpc_security_group_ingress_rule" "https" {
  security_group_id = aws_security_group.serve_vm.id
  cidr_ipv4         = "76.131.61.220/32"
  from_port         = 443
  ip_protocol       = "tcp"
  to_port           = 443
}

resource "aws_vpc_security_group_ingress_rule" "nginx" {
  security_group_id = aws_security_group.serve_vm.id
  cidr_ipv4         = "76.131.61.220/32"
  from_port         = 3000
  ip_protocol       = "tcp"
  to_port           = 3000
}

resource "aws_instance" "serve_app" {
  ami           = "ami-0ce0152a3f6225d58"
  instance_type = "t3.medium"
  associate_public_ip_address = true
  subnet_id = data.aws_subnets.public.ids[0]
  key_name = "aravn-mbp"
  security_groups = [aws_security_group.serve_vm.id]
  root_block_device {
    volume_size = 30
  }

  user_data = <<EOF
#!/bin/bash
## TODO pass in secrets
EOF
  tags = {
    Name = "serve-app"
  }
}