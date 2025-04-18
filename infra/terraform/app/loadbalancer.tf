### SG ###
resource "aws_security_group" "serve_alb" {
  name        = "serve-alb"
  description = "serve alb"
  vpc_id      = data.aws_vpc.this.id
}

resource "aws_vpc_security_group_ingress_rule" "alb_http" {
  security_group_id = aws_security_group.serve_alb.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  to_port           = 80
  ip_protocol       = "tcp"
}

resource "aws_vpc_security_group_ingress_rule" "alb_https" {
  security_group_id = aws_security_group.serve_alb.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
}

### Load Balancer
module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "9.15.0"

  name    = "servev2"
  vpc_id  = data.aws_vpc.this.id
  subnets = data.aws_subnets.public.ids

  # For example only
  enable_deletion_protection = false

  # Security Group
  create_security_group = false
  security_groups       = [aws_security_group.serve_alb.id]


  listeners = {
    http-https-redirect = {
      port     = 80
      protocol = "HTTP"
      redirect = {
        port        = "443"
        protocol    = "HTTPS"
        status_code = "HTTP_301"
      }
    }


    https = {
      port            = 443
      protocol        = "HTTPS"
      ssl_policy      = "ELBSecurityPolicy-TLS13-1-2-2021-06"
      certificate_arn = data.aws_acm_certificate.serve.arn
      additional_certificate_arns = []
      forward = {
        target_group_key = "serve-fe"
      }

      rules = {
        api = {
          priority = 10
          tags = {
            Name = "api"
          }
          actions = [
            {
              type             = "forward"
              target_group_key = "serve-api"
            }
          ]
          conditions = [{
            path_pattern = {
              values = ["/api/*"]
            }
          }]
        }
      }
    }
  }

  target_groups = {
    serve-fe = {
        name = "serve-frontend"
      target_type = "instance"
      target_id   = "i-0c1267556bad83229"
      port        = 3000
      protocol    = "HTTP"
      health_check = {
        enabled             = true
        healthy_threshold   = 5
        interval            = 30
        matcher             = "200"
        path                = "/"
        port                = "traffic-port"
        protocol            = "HTTP"
        timeout             = 5
        unhealthy_threshold = 2
      }
    }
    serve-api = {
        name = "serve-api"
      port      = 8080
      protocol  = "HTTP"
      vpc_id    = data.aws_vpc.this.id
      target_id = "i-0c1267556bad83229"
      health_check = {
        enabled             = true
        healthy_threshold   = 5
        interval            = 30
        matcher             = "404"
        path                = "/"
        port                = "traffic-port"
        protocol            = "HTTP"
        timeout             = 5
        unhealthy_threshold = 2
      }
    }
  }
}