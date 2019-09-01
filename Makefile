.PHONY: build
build:
	go build

.PHONY: commitcheck
commitcheck: clean test

.PHONY: clean
clean:
	rm -f apitest
	go mod tidy

.PHONY: test
test: | depcheck apitest

.PHONY: depcheck
depcheck: 
	#Ignore - working offline, uncomment when back online
	go get -u
	@echo "Did you start the Local Database and other dependencies?"

.PHONY: deploydev
deploydev: test
	rm -f binarypkg.zip
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o binarypkg runner.go
	zip -j binarypkg.zip binarypkg infra/dev/.env
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
