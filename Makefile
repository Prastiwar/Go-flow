.PHONY: default
default: build

all: clean get-deps build test

build:
	mkdir -p bin
	go build -o bin/Go-flow ./...

test: build
	go test -short -coverprofile=bin/cov.out `go list ./... | grep -v vendor/`
	go tool cover -func=bin/cov.out

clean:
	rm -rf ./bin
