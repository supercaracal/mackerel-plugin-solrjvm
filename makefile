SHELL := /bin/bash

all:
	@$(MAKE) --no-print-directory lint
	@$(MAKE) --no-print-directory test
	@$(MAKE) --no-print-directory build

build: mackerel-plugin-solrjvm

mackerel-plugin-solrjvm: main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $@

lint:
	@go vet
	@golint -set_exit_status

test:
	@go test

.PHONY: all build lint test
