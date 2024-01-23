provider "aws" {
  region = "us-east-1"
}

module "onboarding-repo-scan-lambda" {
  count  = 1
  source = "terraform-aws-modules/lambda/aws"

  function_name = "onboarding-repo-scan-lambda"
  description   = "AWS Lambda function to scan repos onboarding"
  handler       = "main"
  timeout       = 30
  memory_size   = 128
  source_path = "./main"

  runtime        = "go1.x"
  package_type   = "Zip"

  environment_variables = {
    QUEUE_URL = "https://sqs.us-east-1.amazonaws.com/586929748635/onboarding-repo-scan"
  }

  event_source_mapping = {
    sqs = {
      service          = "sqs"
      event_source_arn = "arn:aws:sqs:us-east-1:586929748635:onboarding-repo-scan"
      batch_size       = 1
    }
  }

  allowed_triggers = {
    sqs = {
      principal  = "sqs.amazonaws.com"
      source_arn = "arn:aws:sqs:us-east-1:586929748635:onboarding-repo-scan"
    }
  }

  attach_policies    = true
  number_of_policies = 3
  policies = [
    "arn:aws:iam::aws:policy/AmazonSQSFullAccess",
    "arn:aws:iam::aws:policy/SecretsManagerReadWrite",
    "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess"
  ]

  attach_network_policy                   = true
  create_current_version_allowed_triggers = false
}