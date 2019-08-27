.PHONY: all
all: build

BIN_DIR=$(PWD)/bin

install-reflex:
	go get github.com/cespare/reflex

dependencies:
	go mod download

build: dependencies
	go build -o $(BIN_DIR)/mockrift cmd/web/*.go

run: build
	$(BIN_DIR)/mockrift

run-reflex: install-reflex
	reflex -c reflex.conf

help: build
	$(BIN_DIR)/mockrift -help

