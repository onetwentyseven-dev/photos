resource "random_string" "bucket_name" {
  length  = 6
  upper   = false
  special = false
}

resource "aws_s3_bucket" "photos" {
  bucket = "photos-${random_string.bucket_name.result}-${var.region}"
}

resource "aws_s3_bucket" "photos_uploads" {
  bucket = "photos-${random_string.bucket_name.result}-uploads-${var.region}"
}

resource "aws_s3_bucket_acl" "photo_uploads" {
  bucket     = aws_s3_bucket.photos_uploads.id
  acl        = "private"
  depends_on = [aws_s3_bucket_ownership_controls.photo_uploads]
}

resource "aws_s3_bucket_ownership_controls" "photo_uploads" {
  bucket = aws_s3_bucket.photos_uploads.id
  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

resource "aws_s3_bucket_policy" "photo_uploads" {
  bucket = aws_s3_bucket.photos_uploads.id
  policy = data.aws_iam_policy_document.photo_uploads.json
}

data "aws_iam_policy_document" "photo_uploads" {
  statement {
    sid    = "AllowCloudfrontOAI"
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
    ]

    principals {
      type = "AWS"
      identifiers = [
        aws_cloudfront_origin_access_identity.upload.iam_arn
      ]
    }
    resources = [
      "${aws_s3_bucket.photos_uploads.arn}/*"
    ]
  }
}
