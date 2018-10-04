all: build

build:
	go install

test: format lint
	go test -cover -race ./...

format:
	go fmt ./...

lint:
	golint ./...