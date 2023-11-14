resource "aws_apigatewayv2_api" "photos" {
  name          = "photos"
  protocol_type = "HTTP"

  cors_configuration {
    allow_headers = ["*"]
    allow_methods = ["GET", "POST", "PATCH", "DELETE"]
    allow_origins = ["*"]
  }
}

resource "aws_apigatewayv2_domain_name" "photos" {
  domain_name = local.default_domain
  domain_name_configuration {
    certificate_arn = aws_acm_certificate.photos.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_stage" "photos" {
  name        = "$default"
  api_id      = aws_apigatewayv2_api.photos.id
  auto_deploy = true
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.access_logs.arn
    format = jsonencode({
      httpMethod              = "$context.httpMethod"
      ip                      = "$context.identity.sourceIp"
      protocol                = "$context.protocol"
      requestId               = "$context.requestId"
      requestTime             = "$context.requestTime"
      responseLength          = "$context.responseLength"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      integrationErrorMessage = "$context.integrationErrorMessage"
      integrationLatency      = "$context.integrationLatency"
      integrationStatus       = "$context.integrationStatus"
    })
  }
}

resource "aws_apigatewayv2_api_mapping" "photos" {
  api_id      = aws_apigatewayv2_api.photos.id
  domain_name = aws_apigatewayv2_domain_name.photos.id
  stage       = aws_apigatewayv2_stage.photos.id
}
