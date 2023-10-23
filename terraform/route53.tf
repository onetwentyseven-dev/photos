data "cloudflare_zone" "onetwentyseven" {
  name = "onetwentyseven.dev"
}

resource "cloudflare_record" "ns" {
  for_each = toset(aws_route53_zone.photos.name_servers)
  zone_id  = data.cloudflare_zone.onetwentyseven.id
  name     = "photos"
  value    = each.value
  type     = "NS"
  proxied  = false
  ttl      = 300
}

resource "aws_route53_zone" "photos" {
  name = local.default_domain
}

resource "aws_route53_record" "certificate_validation" {
  for_each = {
    for dvo in aws_acm_certificate.photos.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = aws_route53_zone.photos.zone_id
}

output "zone_ns_records" {
  value = aws_route53_zone.photos.name_servers
}

resource "aws_route53_record" "upload" {
  name    = "upload"
  type    = "A"
  zone_id = aws_route53_zone.photos.zone_id

  alias {
    name                   = aws_cloudfront_distribution.upload.domain_name
    zone_id                = aws_cloudfront_distribution.upload.hosted_zone_id
    evaluate_target_health = false
  }

}

resource "aws_route53_record" "photos" {
  name    = aws_apigatewayv2_domain_name.photos.domain_name
  type    = "A"
  zone_id = aws_route53_zone.photos.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.photos.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.photos.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}
