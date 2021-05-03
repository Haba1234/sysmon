BIN := "./bin/sysmon"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/sysmon

run: build
	$(BIN) -config ./configs/config.toml

test:
	go test -v -count=100 -race -timeout=1m ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.39.0

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

CLIENT_BIN := "./bin/client"

client-build:
	go build -v -o $(CLIENT_BIN) ./cmd/client

client-run: client-build
	$(CLIENT_BIN)

.PHONY: build run test lint generate client-build client-run
