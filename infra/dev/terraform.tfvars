#Region defaults to us-east-1

#Shared variables
tags= {
    "TF" = "lifeapp_dev"
}

#lambda memory drives cpu, can bump up to 512 or down to 128 based on performance results later
awsregion="us-east-1"
lambda_memory_size=256
lambda_timeout=15
environment_variables= {
    "LAMBDA_ENV" = "true"
}

#DB variables - note:all settings here are shared except for table name, would split up into separate variables in higher environment
dyamodb_name="category_dev"
dyamodb_read_capacity="1"  
dyamodb_write_capacity="1"  
dyamodb_hash_key="CategoryHashKey"
dyamodb_range_key="CategorySortKey"
dyamodb_stream_enabled="false"
dyamodb_stream_view_type=""
dynamodb_table_attributes=[
    {
        name = "CategoryHashKey",
        type = "S",
    },
    {
        name = "CategorySortKey",
        type = "S",
    }
] 
