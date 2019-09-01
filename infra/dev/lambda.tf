
#First upload the file -- assumes the bucket has been created centrally / externally
resource "aws_s3_bucket_object" "binary" {
  bucket = "lifeapp-terraform-dev"
  key    = "dev-lifefapp.zip"

  #source is a local path
  source = var.binaryfilepath
  etag   = md5(filebase64(var.binaryfilepath))
}

# lambda_s3 is used as direct uploads have a smaller file limit
# To force reload example --> terraform taint aws_lambda_function.process_lambda
resource "aws_lambda_function" "process_lambda" {
  function_name = "LifeApp"

  # The bucket name as created earlier with "aws s3api create-bucket"
  s3_bucket = aws_s3_bucket_object.binary.bucket
  s3_key    = aws_s3_bucket_object.binary.key

  #Versioning to be done later
  #s3_key    = "v${var.app_version}/example.zip"

  #The executable file name value. For example, "myHandler" would call the main function in the package “main” of the myHandler executable program.
  handler = "binarypkg"
  runtime = "go1.x"

  memory_size = var.lambda_memory_size
  role        = aws_iam_role.demo_lambda_exec.arn
  timeout     = var.lambda_timeout

  environment {
    variables = var.environment_variables
  }

  tags = var.tags
}

# IAM role which dictates what other AWS services the Lambda function
# may access.

resource "aws_iam_role" "demo_lambda_exec" {
  name = "lambda_demo_role"

  assume_role_policy = file("lambda_exec_iam.json")
}


#This exists in my env already
resource "aws_iam_policy" "cloudwatch-demo-policy" {
  name        = "demo-cloudwatch-policy"
  description = "grants access to log in cloudwatch"
  policy      = file("lambda_cloudwatch_iam.json")
}

resource "aws_iam_role_policy_attachment" "lambda-cloudwatch-policy" {
  role       = aws_iam_role.demo_lambda_exec.name
  policy_arn = aws_iam_policy.cloudwatch-demo-policy.arn
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.process_lambda.arn
  principal     = "apigateway.amazonaws.com"

  # The /*/* portion grants access from any method on any resource
  # within the API Gateway "REST API".
  source_arn = "${data.terraform_remote_state.core_infra_dev.outputs.api_gw.execution_arn}/*/*/lifeapp/*"
}

