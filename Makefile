SHELL := /bin/bash

default: build

build:
	go build -v ./...

test: build
	go test -v -coverprofile=cov.out `go list ./... | grep -v vendor/`
