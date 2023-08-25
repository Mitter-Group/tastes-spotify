resource "aws_dynamodb_table" "user_tastes" {
  name           = var.dynamo_user_tastes_table_name
  billing_mode   = "PAY_PER_REQUEST"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  hash_key  = "PK"
  range_key = "SK"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

resource "aws_dynamodb_table" "tastes_tags" {
  name           = var.dynamo_tastes_tags_table_name
  billing_mode   = "PAY_PER_REQUEST"

  attribute {
    name = "PK"
    type = "S"
  }

  hash_key = "PK"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}
