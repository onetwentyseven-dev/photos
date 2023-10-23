resource "aws_cloudwatch_log_group" "access_logs" {
  name              = "/aws/apigateway/${aws_apigatewayv2_api.photos.name}"
  retention_in_days = 3
}
