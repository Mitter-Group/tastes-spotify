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

variable "dynamo_user_tastes_table_name" {
  description = "El nombre de la tabla DynamoDB para los gustos del usuario."
  type        = string
}

variable "dynamo_tastes_tags_table_name" {
  description = "El nombre de la tabla DynamoDB para las etiquetas de gustos."
  type        = string
}
