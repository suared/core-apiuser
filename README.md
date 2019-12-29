Example API app that uses the core library.  

This example uses the new repository hooks as well to simplify Dynamo access.  To start, after getting this repo:

1.  Do a "find-replace all" - github.com/suared/core-apiuser ==> yourpackage
1.1.  Remove the .git directory or copy everything except the .git directory into your target repo  (if the former, create a new repo)
2.  Starting from the model and working up, replace each layer with your intended package.  Generally rename the files to target and then update using the same overall style
3.  run "make test" from your new base dir. This is your test target and will run your test files in order from the model tier upward
4.  update your infra/ to be relevant to you, especially the .env files in the base and dev/ which will be used when deploying to dev as well as the variable files
5.  do the same above to each respective environment you will have

Done

Critical notes:
* Test in order to minimize confusion - packages from bottom (model) to top (api).  
* Do the same as you deploy to new environments to validate each step