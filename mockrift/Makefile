.PHONY: all
all: build

BIN_DIR=$(PWD)/bin

install-reflex:
	go get -u github.com/cespare/reflex

install-delve:
	go get -u github.com/go-delve/delve/cmd/dlv

dependencies:
	go mod download

build: dependencies
	go build -o $(BIN_DIR)/mockrift cmd/web/*.go

run: build
	$(BIN_DIR)/mockrift

run-reflex:
	reflex -c reflex.conf

run-delve: dependencies
	dlv debug cmd/web/*.go --headless=true --listen=":40000" --log --accept-multiclient --api-version=2 --output=./bin/__debug_bin

help: build
	$(BIN_DIR)/mockrift -help

