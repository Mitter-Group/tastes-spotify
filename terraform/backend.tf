terraform {
  backend "s3" {
    bucket = var.backend_state_bucket
    key    = var.backend_state_key
    region = var.aws_region
  }
}