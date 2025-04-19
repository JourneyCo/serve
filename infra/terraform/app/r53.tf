resource "aws_route53_record" "serve" {
  provider = aws.dns
  zone_id  = data.aws_route53_zone.this.zone_id
  name     = "serve.ravn.systems"
  type     = "A"
  alias {
    name                   = module.alb.dns_name
    zone_id                = module.alb.zone_id
    evaluate_target_health = false
  }
}