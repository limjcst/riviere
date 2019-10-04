SOURCES=$(shell go list ./...)

all: build

build:
	go build
	go generate

test: format
	go test -cover -race -v -coverprofile=coverage.out ./...

format:
	go fmt $(SOURCES)

lint:
	golint $(SOURCES)
