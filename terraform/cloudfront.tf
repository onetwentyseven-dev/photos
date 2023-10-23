locals {
  uploads_origin_id = "s3-uploads-${random_string.origin_id.result}"
}

resource "random_string" "origin_id" {
  length  = 6
  upper   = false
  special = false
}

data "aws_cloudfront_cache_policy" "managed_caching_disabled" {
  name = "Managed-CachingDisabled"
}

data "aws_cloudfront_origin_request_policy" "managed_all_viewer" {
  name = "Managed-AllViewerExceptHostHeader"
}

resource "aws_cloudfront_origin_access_identity" "upload" {
  comment = "Uploads"
}

resource "aws_cloudfront_response_headers_policy" "upload" {
  name = "upload-response-policy"

  cors_config {
    access_control_allow_credentials = true

    access_control_allow_origins {
      items = ["https://photos.onetwentyseven.dev"]
    }
    access_control_allow_methods {
      items = ["OPTIONS", "PUT"]
    }
    access_control_allow_headers {
      items = ["Content-Type"]
    }

    origin_override = true
  }
}

resource "aws_cloudfront_distribution" "upload" {
  enabled         = true
  is_ipv6_enabled = true

  origin {
    domain_name = aws_s3_bucket.photos_uploads.bucket_regional_domain_name
    origin_id   = local.uploads_origin_id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.upload.cloudfront_access_identity_path
    }
  }

  aliases = [
    "upload.${local.default_domain}",
  ]

  default_cache_behavior {
    target_origin_id           = local.uploads_origin_id
    cache_policy_id            = data.aws_cloudfront_cache_policy.managed_caching_disabled.id
    allowed_methods            = ["HEAD", "DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"]
    cached_methods             = ["HEAD", "GET", "OPTIONS"]
    viewer_protocol_policy     = "redirect-to-https"
    response_headers_policy_id = aws_cloudfront_response_headers_policy.upload.id
    origin_request_policy_id   = data.aws_cloudfront_origin_request_policy.managed_all_viewer.id

    lambda_function_association {
      event_type = "viewer-request"
      lambda_arn = module.photos_edge_validation.function_qualified_arn
    }
    lambda_function_association {
      event_type = "origin-request"
      lambda_arn = module.photos_edge_validation.function_qualified_arn
    }
    lambda_function_association {
      event_type = "origin-response"
      lambda_arn = module.photos_edge_validation.function_qualified_arn
    }

    lambda_function_association {
      event_type = "viewer-response"
      lambda_arn = module.photos_edge_validation.function_qualified_arn
    }

  }



  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate.photos.arn
    minimum_protocol_version = "TLSv1.2_2021"
    ssl_support_method       = "sni-only"
  }
}
