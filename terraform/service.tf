# resource "aws_s3_bucket" "beanstalk_app_bucket" {
#   bucket = var.s3_bucket_name

#   tags = {
#     Name        = "Beanstalk App Bucket"
#     Environment = "Dev"
#   }
# }

# resource "aws_elastic_beanstalk_application" "myapp" {
#   name        = var.eb_app_name
#   description = "Spotify Service"
# }

# resource "aws_elastic_beanstalk_environment" "myenv" {
#   name                = var.eb_env_name
#   application         = aws_elastic_beanstalk_application.myapp.name
#   solution_stack_name = var.eb_solution_stack
#   version_label       = var.eb_version_label

#   setting {
#     namespace = "aws:elasticbeanstalk:environment"
#     name      = "EnvironmentType"
#     value     = "LoadBalanced"
#   }

#   setting {
#     namespace = "aws:autoscaling:launchconfiguration"
#     name      = "IamInstanceProfile"
#     value     = "ElasticBeanstalkInstanceProfile"
#   }
# }

# resource "aws_elastic_beanstalk_application_version" "myapp_version" {
#   name        = var.eb_version_label
#   application = aws_elastic_beanstalk_application.myapp.name

#   bucket  = var.s3_bucket_name
#   key     = "spotify/spotify.zip"
# }
