GO_BUILD= go build
GOFLAGS= GOOS=linux CGO_ENABLED=0

build:
	$(GOFLAGS) $(GO_BUILD) -tags netgo -a -v -installsuffix cgo -o bin/generic-bot main.go
.PHONY: go-build

vendor: go.mod
	go mod tidy
.PHONY: vendor

all: vendor build
.PHONY: all

default: all
.PHONY: default