resource "aws_security_group" "serve_vm" {
  name        = "serve-vm"
  description = "serve vm SG"
  vpc_id      = data.aws_vpc.this.id
}

resource "aws_vpc_security_group_ingress_rule" "allow_lb" {
  security_group_id = aws_security_group.serve_vm.id
  cidr_ipv4         = data.aws_vpc.this.cidr_block
  ip_protocol       = "-1"
}

# resource "aws_vpc_security_group_ingress_rule" "ssh" {
#   security_group_id = aws_security_group.serve_vm.id
#   cidr_ipv4         = "76.131.61.220/32"
#   from_port         = 22
#   ip_protocol       = "tcp"
#   to_port           = 22
# }

resource "aws_vpc_security_group_egress_rule" "egress" {
  security_group_id = aws_security_group.serve_vm.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

resource "aws_instance" "serve_app" {
  ami                         = var.ami_id
  instance_type               = var.instance_type
  associate_public_ip_address = true
  subnet_id                   = data.aws_subnets.public.ids[0]
  key_name                    = var.key_name
  vpc_security_group_ids      = [aws_security_group.serve_vm.id]
  root_block_device {
    volume_size = var.root_block_size
  }

  user_data = templatefile("templates/userdata.tpl", {
    dev_mode              = var.dev_mode,
    serve_day             = var.serve_day,
    api_port              = var.api_port,
    db_host               = var.db_host,
    db_port               = var.db_port,
    db_user               = var.db_user,
    db_pass               = var.db_pass,
    db_name               = var.db_name,
    auth0_domain          = var.auth0_domain,
    auth0_audience        = var.auth0_audience,
    auth0_client_id       = var.auth0_client_id,
    auth0_client_secret   = var.auth0_client_secret,
    mailtrap_host         = var.mailtrap_host,
    mailtrap_key          = var.mailtrap_key,
    mailtrap_from         = var.mailtrap_from,
    mailtrap_port         = var.mailtrap_port,
    mailtrap_user         = var.mailtrap_user,
    mailtrap_pass         = var.mailtrap_pass,
    clearstream_api_key   = var.clearstream_api_key,
    clearstream_text_from = var.clearstream_text_from,
    google_key            = var.google_key
  })
  tags = {
    Name = "serve-app"
  }
}