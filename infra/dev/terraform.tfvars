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

