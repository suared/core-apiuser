###AWS Config Section ###
variable "awsregion" {
}

variable "awscredentialsfile" {
}

variable "awsprofile" {
}

#Shared
variable "tags" {
  type = map(string)
}

variable "core_infra_state_filepath" {
  
}

#####Lambda Setup Section ####
#Versioning to be done later
#variable "app_version" {}
variable "lambda_memory_size" {
}

variable "lambda_timeout" {
}

variable "environment_variables" {
  type = map(string)
}

variable "binaryfilepath" {
}

#####API GW Setup Section ####

