resource "random_string" "bucket_name" {
  length  = 6
  upper   = false
  special = false
}

resource "aws_s3_bucket" "photos" {
  bucket = "photos-${random_string.bucket_name.result}-${var.region}"
}


# resource "aws_s3_bucket_acl" "photos" {
#   bucket     = aws_s3_bucket.photos.id
#   acl        = "private"
#   depends_on = [aws_s3_bucket_ownership_controls.photos]
# }

# resource "aws_s3_bucket_ownership_controls" "photos" {
#   bucket = aws_s3_bucket.photos.id
#   rule {
#     object_ownership = "BucketOwnerEnforced"
#   }
# }

resource "aws_s3_bucket_policy" "photos" {
  bucket = aws_s3_bucket.photos.id
  policy = data.aws_iam_policy_document.photos.json
}

data "aws_iam_policy_document" "photos" {
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
      "${aws_s3_bucket.photos.arn}/*"
    ]
  }
}


resource "aws_s3_bucket_notification" "photo_uploads" {
  bucket      = aws_s3_bucket.photos_uploads.id
  eventbridge = true
}

resource "aws_cloudwatch_event_rule" "process_photo_upload" {
  name        = "process-photo-upload"
  description = "Triggers a lambda to pull a file out of S3, generate a thumbnail and extract the exif data and save it to the database"

  event_pattern = jsonencode({
    source      = ["aws.s3"],
    detail-type = ["Object Created"],
    detail = {
      bucket = {
        name = [aws_s3_bucket.photos.id]
      },
      object = {
        key = [
          "originals/*"
        ]
      }
    }

  })
}

resource "aws_cloudwatch_event_target" "process_photo_upload" {
  rule = aws_cloudwatch_event_rule.process_photo_upload.name
  arn  = module.photos_processor.function_arn
}


resource "aws_lambda_permission" "allow_eventbridge" {
  action        = "lambda:InvokeFunction"
  function_name = module.photos_processor.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.process_photo_upload.arn
}
