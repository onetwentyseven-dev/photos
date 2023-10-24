locals {
  lambda_basic_policy = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  lambda_vpc_policy   = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
  lambda_policy       = var.vpc_id != "" ? local.lambda_vpc_policy : local.lambda_basic_policy
  parameter_store_arn = var.paramstore_prefix != null ? "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:parameter${var.paramstore_prefix}/*" : null
  lambda_handler      = var.function_package_location != "" ? var.function_package_location : "${path.module}/assets/lambda_base.zip"
}

data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

resource "aws_lambda_function" "function" {
  function_name    = var.function_name
  filename         = local.lambda_handler
  source_code_hash = filebase64sha256(local.lambda_handler)
  runtime          = var.function_runtime
  role             = aws_iam_role.lambda.arn
  handler          = var.function_handler
  timeout          = var.function_timeout
  memory_size      = var.function_memory
  architectures    = ["x86_64"]
  publish          = var.enable_versions
}

data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = [
        "edgelambda.amazonaws.com",
        "lambda.amazonaws.com"
      ]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "${var.function_name}-lambda-${data.aws_region.current.name}"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_role_policy_attachment" "lambda" {
  role       = aws_iam_role.lambda.name
  policy_arn = local.lambda_policy

  lifecycle {
    create_before_destroy = true
  }
}

# resource "aws_iam_role_policy" "lambda_paramstore" {
#   count  = local.parameter_store_arn != null ? 1 : 0
#   name   = "paramstore-access"
#   role   = aws_iam_role.lambda.name
#   policy = data.aws_iam_policy_document.lambda_paramstore[0].json
# }

# data "aws_iam_policy_document" "lambda_paramstore" {
#   count = local.parameter_store_arn != null ? 1 : 0
#   statement {
#     effect = "Allow"
#     actions = [
#       "ssm:GetParameter",
#       "ssm:GetParameters",
#       "ssm:GetParametersByPath",
#       "ssm:DescribeParameters",
#     ]
#     resources = [
#       local.parameter_store_arn
#     ]
#   }
# }


# resource "aws_iam_role_policy" "lambda_cloudwatch" {
#   name   = "cloudwatch-access"
#   role   = aws_iam_role.lambda.name
#   policy = data.aws_iam_policy_document.lambda_cloudwatch.json
# }

# data "aws_iam_policy_document" "lambda_cloudwatch" {
#   statement {
#     effect = "Allow"
#     actions = [
#       "logs:CreateLogStream",
#       "logs:PutLogEvents"
#     ]
#     resources = [
#       aws_cloudwatch_log_group.lambda.arn
#     ]
#   }
# }

resource "aws_iam_role_policy" "lambda_custom" {
  for_each = var.additional_role_policies
  name     = each.key
  role     = aws_iam_role.lambda.name
  policy   = each.value
}

# resource "aws_cloudwatch_log_group" "lambda" {
#   name              = "/aws/lambda/${var.function_name}"
#   retention_in_days = var.log_retention_in_days
# }

# resource "aws_security_group" "lambda" {
#   count       = var.vpc_id != "" ? 1 : 0
#   name        = "${var.function_name}-lambda-egress"
#   description = "Allow lambda egress"
#   vpc_id      = var.vpc_id
# }

# resource "aws_security_group_rule" "egress" {
#   count             = var.vpc_id != "" ? 1 : 0
#   type              = "egress"
#   from_port         = 0
#   to_port           = 0
#   protocol          = "-1"
#   security_group_id = aws_security_group.lambda[0].id

#   cidr_blocks = [
#     "0.0.0.0/0",
#   ]
# }
