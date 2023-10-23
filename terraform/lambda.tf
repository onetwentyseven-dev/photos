locals {
  default_envs = {
    "APP_URL"             = "https://${aws_apigatewayv2_domain_name.photos.domain_name}",
    "AUTH0_DOMAIN"        = "onetwentyseven.us.auth0.com",
    "AUTH0_CLIENT_ID"     = "qXN9bddpwY7EpoF7RTYno2Y3GlkGV5tA"
    "AUTH0_CALLBACK_PATH" = "/oauth/callback"
    "DB_HOST"             = "aws.connect.psdb.cloud"
    "DB_NAME"             = "photos"
    "DB_USER"             = "2p2p2trqlegwtwae66q2"
    "MODE"                = "lambda"
    "ENVIRONMENT"         = "production"


  }
}


module "photos_handler" {
  source  = "onetwentyseven-dev/lambda/aws"
  version = "~> 1"

  function_name     = "photos-handler"
  function_runtime  = "provided.al2"
  function_memory   = 128
  paramstore_prefix = "/photos"

  log_retention_in_days = 7

  environment_variables = local.default_envs
}



module "photos_edge_validation" {
  source = "./modules/photos_edge_validation"

  function_name             = "photos-edge-validation"
  function_runtime          = "nodejs18.x"
  function_memory           = 128
  function_timeout          = 5
  function_handler          = "index.handler"
  enable_versions           = true
  function_package_location = "${path.module}/assets/photos-edge-validation.zip"

}

module "photos_routes" {
  source  = "onetwentyseven-dev/lambda-apigw/aws"
  version = "~> 1"

  apigw_id      = aws_apigatewayv2_api.photos.id
  function_name = module.photos_handler.function_name

  routes = [
    "GET /static/{proxy+}",

    "GET /api/auth/validate",
    "POST /api/image/metadata",

    "GET /",
    "GET /login",
    "GET /logout",
    "GET /oauth/callback",

    "GET /dashboard",
    "GET /dashboard/upload",
    "POST /dashboard/upload",
  ]

}

module "photos_processor" {
  source  = "onetwentyseven-dev/lambda/aws"
  version = "~> 1"

  function_name    = "photos-processor"
  function_runtime = "provided.al2"
  function_memory  = 128

  log_retention_in_days = 7

  environment_variables = local.default_envs
}

module "photos_migrations" {
  source  = "onetwentyseven-dev/lambda/aws"
  version = "~> 1"

  function_name     = "photos-migrations"
  function_runtime  = "provided.al2"
  function_memory   = 128
  paramstore_prefix = "/photos"

  log_retention_in_days = 7

  environment_variables = local.default_envs
}
