resource "aws_route53_record" "serve" {
  count    = var.create_dns_record ? 1 : 0
  provider = aws.dns
  zone_id  = data.aws_route53_zone.this[0].zone_id
  name     = "serve.${var.domain}"
  type     = "A"
  alias {
    name                   = module.alb.dns_name
    zone_id                = module.alb.zone_id
    evaluate_target_health = false
  }
}