APP?=banking_app_authServer_authServer

.PHONY: build
## build: build the application
build: clean
	@echo "Building..."
	@go build -o ${APP} main.go

.PHONY: run
## run: runs the go run build-binary
run: build
	./banking_app_authServer

.PHONY: watch
## watch: watch the project for go file changes
watch:
	ulimit -n 1000 
	reflex -s -r '\.go$$' make run

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning..."
	@go clean

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...


.PHONY: help
## help: prints this help message
help:

	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'