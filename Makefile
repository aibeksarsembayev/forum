.PHONY: build
build:
	go build -v ./cmd/web

.DEFAULT_GOAL := build