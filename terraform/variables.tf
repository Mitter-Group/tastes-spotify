variable "aws_region" {
  description = "La región de AWS para los recursos."
  type        = string
  default     = "us-west-1"
}

variable "aws_profile" {
  description = "El perfil de AWS a utilizar."
  type        = string
  default     = "personal"
}

variable "s3_bucket_name" {
  description = "El nombre del bucket S3 para el servicio Beanstalk."
  type        = string
}

variable "eb_app_name" {
  description = "El nombre de la aplicación Elastic Beanstalk."
  type        = string
}

variable "eb_env_name" {
  description = "El nombre del entorno Elastic Beanstalk."
  type        = string
}

variable "eb_solution_stack" {
  description = "La pila de soluciones para Elastic Beanstalk."
  type        = string
}

variable "eb_version_label" {
  description = "Etiqueta de versión de la aplicación Elastic Beanstalk."
  type        = string
}

variable "table_name" {
  description = "Name of the DynamoDB table"
  type        = string
}

variable "user_data_table_name" {
  description = "Name of the DynamoDB table"
  type        = string
}

variable "billing_mode" {
  description = "Billing mode for the DynamoDB table"
  type        = string
  default     = "PAY_PER_REQUEST"
}

variable "hash_key" {
  description = "Hash key for the DynamoDB table"
  type        = string
}

variable "range_key" {
  description = "Range key for the DynamoDB table"
  type        = string
}

variable "backend_state_bucket" {
  description = "Terraform state bucket"
  type        = string
}

variable "backend_state_key" {
  description = "Teraform state file key"
  type        = string
}
