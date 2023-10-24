resource "aws_iam_role_policy" "photos_processing" {
  name   = "photo-processing-read-objects"
  role   = module.photos_processor.iam_role_name
  policy = data.aws_iam_policy_document.photos_processing.json
}

data "aws_iam_policy_document" "photos_processing" {
  statement {
    # Allow the lambda to read from the S3 bucket aws_s3_bucket.photo_uploads.name
    sid    = "AllowLambdaToReadFromS3"
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:DeleteObject",
    ]
    resources = [
      "${aws_s3_bucket.photos_uploads.arn}/*"
    ]
  }
  statement {
    # Allow the lambda to write to the S3 bucket aws_s3_bucket.photos.name
    sid    = "AllowLambdaToWriteToS3"
    effect = "Allow"
    actions = [
      "s3:PutObject",
    ]
    resources = [
      "${aws_s3_bucket.photos.arn}/*"
    ]
  }
}
