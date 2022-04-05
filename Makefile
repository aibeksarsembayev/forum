.PHONY: build

build:
	go build -v ./cmd/web

.DEFAULT_GOAL := build

run:
	docker run --name forumapp -p 5050:4000 forum 