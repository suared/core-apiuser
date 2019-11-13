.PHONY: build
build:
	go build

.PHONY: commitcheck
commitcheck: clean test

.PHONY: clean
clean:
	rm -f apitest
	rm -f servicetest
	rm -f repositorytest
	rm -f modeltest
	go mod tidy

.PHONY: test
test: | depcheck modeltest repositorytest servicetest apitest

.PHONY: depcheck
depcheck: 
	#Ignore - working offline, uncomment when back online
	#go get -u #Not required with go mod
	@echo "Did you start the Local Database and other dependencies?"

.PHONY: deploylocal
deploylocal: test
	go run runner.go

.PHONY: deploydev
deploydev: test
	rm -f binarypkg.zip
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o binarypkg runner.go
	zip -j binarypkg.zip binarypkg infra/dev/.env
	#First time there is no resource to taint - terraform init also needs to be run first to setup as one time tasks
	#cd infra/dev && terraform apply
	cd infra/dev && terraform taint aws_lambda_function.process_lambda && terraform apply

.PHONY: deploydeve2e
deploydeve2e:
	PROCESS_ENV=test && export PROCESS_ENV && cd api && go test


#TODO: Change most of these to use stamps naming convention/ hidden files when fully tested and functional
apitest: api api/*.go
	cd api && go test -cover -coverprofile=coverage.out
	cd api && go vet
	cd api && go tool cover -html=coverage.out
	@touch $@

servicetest: service service/*.go
	cd service && go test -cover -coverprofile=coverage.out
	cd service && go vet
	cd service && go tool cover -html=coverage.out
	@touch $@

repositorytest: repository repository/*.go
	cd repository && go test -cover -coverprofile=coverage.out
	cd repository && go vet
	cd repository && go tool cover -html=coverage.out
	@touch $@

modeltest: model model/*.go
	cd model && go test -cover -coverprofile=coverage.out
	cd model && go vet
	cd model && go tool cover -html=coverage.out
	@touch $@
