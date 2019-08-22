.PHONY: all
all: build

BIN_DIR=$(PWD)/bin

dependencies:
	go mod download

build: dependencies
	go build -o $(BIN_DIR)/graphql cmd/web/*.go

run: build
	$(BIN_DIR)/graphql

help: build
	$(BIN_DIR)/graphql -help

install-reflex:
	go get github.com/cespare/reflex

debug: dependencies install-reflex
	reflex -c reflex.conf
