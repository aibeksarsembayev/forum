.PHONY: build
build:
	go build -v ./cmd/web

.PHONY: run

run:
	docker run -p 3000:4000 -v forum:/app --rm --name forum forum:volumes

stop:
	docker stop forum	
	
.DEFAULT_GOAL := build