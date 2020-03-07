provider "aws" {
  region                  = var.awsregion
  shared_credentials_file = var.awscredentialsfile
  profile                 = var.awsprofile
}


data "terraform_remote_state" "core_infra_dev" {
  backend = "local"

  config = {
    path = "${var.core_infra_state_filepath}"
  }
}

resource "aws_api_gateway_resource" "mservice" {
  rest_api_id = data.terraform_remote_state.core_infra_dev.outputs.api_gw.id
  parent_id   = data.terraform_remote_state.core_infra_dev.outputs.api_gw.root_resource_id
  path_part   = "lifeapp"
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = data.terraform_remote_state.core_infra_dev.outputs.api_gw.id
  parent_id   = aws_api_gateway_resource.mservice.id
  path_part   = "{proxy+}"
}

#The special path_part value "{proxy+}" activates proxy behavior, which means that this resource will match any request path. 
resource "aws_api_gateway_method" "proxy" {
  rest_api_id   = aws_api_gateway_resource.proxy.rest_api_id
  resource_id   = aws_api_gateway_resource.proxy.id
  http_method   = "ANY"
  authorization = "NONE"
  #authorization = data.terraform_remote_state.core_infra_dev.outputs.api_auth.type
  #authorizer_id = data.terraform_remote_state.core_infra_dev.outputs.api_auth.id

  request_parameters = {
    "method.request.path.proxy" = true
  }
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = aws_api_gateway_resource.proxy.rest_api_id
  resource_id = aws_api_gateway_method.proxy.resource_id
  http_method = aws_api_gateway_method.proxy.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.awsregion}:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:179984542401:function:$${stageVariables.lifeapp_fx}/invocations"
  uri                     = aws_lambda_function.process_lambda.invoke_arn
}

resource "aws_api_gateway_deployment" "process_gw" {
  depends_on = [
    aws_api_gateway_integration.lambda,
  ]

  rest_api_id = aws_api_gateway_resource.proxy.rest_api_id
  stage_name  = "test"
  variables = {
    "lifeapp_fx" = aws_lambda_function.process_lambda.function_name
  }
}

