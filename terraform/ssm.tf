resource "aws_ssm_parameter" "auth0_client_secret" {
  name  = "${local.ssm_prefix}/auth0_client_secret"
  value = "ChangeMe123"
  type  = "SecureString"
  lifecycle {
    ignore_changes = [value]
  }
}

resource "random_string" "session_key" {
  special = false
  length  = 64
}

resource "aws_ssm_parameter" "session_key" {
  name  = "${local.ssm_prefix}/session_key"
  value = random_string.session_key.result
  type  = "SecureString"
}

resource "aws_ssm_parameter" "photos_bucket" {
  name  = "${local.ssm_prefix}/photos_bucket"
  value = aws_s3_bucket.photos.bucket
  type  = "String"
}

resource "aws_ssm_parameter" "photo_uploads_bucket" {
  name  = "${local.ssm_prefix}/photos_uploads_bucket"
  value = aws_s3_bucket.photos_uploads.bucket
  type  = "String"
}

resource "aws_ssm_parameter" "db_pass" {
  name  = "${local.ssm_prefix}/db_pass"
  value = "ChangeMe"
  type  = "SecureString"
  lifecycle {
    ignore_changes = [value]
  }
}
