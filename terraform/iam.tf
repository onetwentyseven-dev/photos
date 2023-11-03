data "aws_iam_policy_document" "photos_s3_permissions" {
  statement {
    # Allow the lambda to read from the S3 bucket aws_s3_bucket.photo_uploads.name
    sid    = "AllowLambdaToReadFromS3"
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
    ]
    resources = [
      "${aws_s3_bucket.photos.arn}/*"
    ]
  }
}
