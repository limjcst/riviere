all: build

build:
	go install

test:
	go test -cover -race ./...

lint:
	golint ./...