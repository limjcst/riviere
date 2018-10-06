SOURCES=$(shell go list ./...)

all: build

build:
	dep ensure
	go get github.com/go-swagger/go-swagger/cmd/swagger
	go build
	swagger generate spec -o swagger.json

test: format lint
	go test -cover -race ./...

format:
	go fmt $(SOURCES)

lint:
	golint $(SOURCES)
