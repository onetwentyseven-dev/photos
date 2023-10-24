terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "=5.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.5.1"
    }
  }
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      Repository = "https://github.com/ddouglas/photos"
    }
  }
}


locals {
  ssm_prefix     = "/photos"
  default_domain = "photos.onetwentyseven.dev"
}

variable "region" {
  default = "us-east-1"
  type    = string
}
