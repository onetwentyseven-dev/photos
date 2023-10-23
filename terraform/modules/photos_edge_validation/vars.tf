variable "function_name" {
  type        = string
  description = "The name to give the function"
}

variable "function_runtime" {
  type        = string
  description = "Identifier of the function's runtime"
}

variable "function_memory" {
  type        = number
  description = "Amount of memory to allocate to the function"
}

variable "paramstore_prefix" {
  type        = string
  default     = null
  description = "Prefix to allow ssm:GetParameterByPath action to. Alternatively, use the iam_role_name output to attach your own policies"
}

variable "log_retention_in_days" {
  type        = number
  default     = 7
  description = "Number of days the events in the log streams will be retained before being purged. (Default: 7)"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID to attach security groups to"
  default     = ""
}

variable "vpc_subnet_ids" {
  type        = list(string)
  description = "List of Subnet IDs to associate with the lambda"
  default     = []
}

variable "vpc_security_group_ids" {
  type        = list(string)
  description = "List of Security Group IDs to attach to the lambda function"
  default     = []
}

variable "additional_role_policies" {
  type        = map(string)
  default     = {}
  description = "Additional Policies that should be attached to the IAM Execution Role created for this lambda"
}

variable "environment_variables" {
  type        = map(string)
  description = "list of environment variables to inject into the lambda"
  default     = {}
}

variable "enable_versions" {
  type        = bool
  default     = false
  description = "Publish Lambda Versions for use with AWS Cloudfront Lambda@Edge"
}

variable "function_timeout" {
  type        = number
  default     = 30
  description = "value in seconds at which the lambda times out"
}

variable "function_package_location" {
  default     = ""
  type        = string
  description = "location relative to the terraform project that this lambda's function code is located at. File path should end in .zip"
  validation {
    condition     = var.function_package_location == "" ? true : can(regex(".*\\.zip$", var.function_package_location))
    error_message = "function_package_location must end in .zip"
  }
}

variable "function_handler" {
  default = ""
  type    = string
}
