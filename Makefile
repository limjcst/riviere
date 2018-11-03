SOURCES=$(shell go list ./...)

all: build

env-check:
	dep version
ifeq ($(GOPATH), )
		echo "GOPATH is required"
		exit 1
endif

build: env-check
	dep ensure
	go build
	go generate

test: format
	go test -cover -race -v -coverprofile=coverage.out ./...

format:
	go fmt $(SOURCES)

lint:
	golint $(SOURCES)
