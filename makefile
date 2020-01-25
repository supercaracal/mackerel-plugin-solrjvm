SHELL      := /bin/bash
owner_id   := supercaracal
app_name   := mackerel-plugin-solrjvm
latest_tag := $(shell git describe --abbrev=0 --tags)

all: build test lint

build: mackerel-plugin-solrjvm

mackerel-plugin-solrjvm: main.go
	go build -ldflags="-s -w" -trimpath -o $@

test:
	go test

lint:
	go vet
	golint -set_exit_status

clean:
	@rm -f mackerel-plugin-solrjvm main

cross-compile:
	goxz -d dist/${latest_tag} -z -os windows,darwin,linux -arch amd64,386

upload-assets:
	ghr -u ${owner} -r ${app_name} ${latest_tag} dist/${latest_tag}

.PHONY: all build test lint clean cross-compile upload-assets
