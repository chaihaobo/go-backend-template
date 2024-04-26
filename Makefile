SHELL:=/bin/bash
OS=linux
BINARY=betemplate
BUILD_PREFIX=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

## build all executable files
build:
	@go build -ldflags '-w -s' -o ${BINARY} .

static:
	@${BUILD_PREFIX} go build -ldflags '-w -s' -o ${BINARY} .

clean:
	@rm -f ${BINARY}*

## run this application
run:
	@go run main.go

## generate mock files
generate-mock:
	@mockery --all --keeptree --inpackage


.PHONY: help
## show help
help:
	@echo ''
	@echo 'Usage: be-template build'
	@echo ' make target'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^## (.*)/); \
	if (helpMessage) { \
	helpCommand = substr($$1, 0, index($$1, ":")-1); \
	helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
	printf " %-$(TARGET_MAX_CHAR_NUM)s\t%s\n", helpCommand, helpMessage; \
	} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
