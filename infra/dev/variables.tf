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


#####DB  Setup Section ####
variable "dyamodb_name" {
}


variable "dyamodb_read_capacity" {
}

variable "dyamodb_write_capacity" {
}

variable "dyamodb_hash_key" {
}

variable "dyamodb_range_key" {
}

variable "dyamodb_stream_enabled" {
}

variable "dyamodb_stream_view_type" {
}

variable "dynamodb_table_attributes" {
  default = [
    {
        name = "SampleHashKey",
        type = "S",
    },
    {
        name = "SampleSortKey",
        type = "S",
    }
] 
}
