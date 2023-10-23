output "function_arn" {
  value = aws_lambda_function.function.arn
}

output "function_name" {
  value = aws_lambda_function.function.function_name
}

output "function_invoke_arn" {
  value = aws_lambda_function.function.invoke_arn
}

output "function_qualified_arn" {
  value = var.enable_versions ? aws_lambda_function.function.qualified_arn : null
}
