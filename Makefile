SOURCES=$(shell go list ./...)

all: build

build:
	dep ensure
	go get github.com/go-swagger/go-swagger/cmd/swagger
	go build
	go generate

test: format lint
	go test -cover -race ./...

format:
	go fmt $(SOURCES)

lint:
	golint $(SOURCES)
