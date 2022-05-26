build:
	docker build -t iforum  .

buildmulti:
	docker build -t iforum:multistage -f Dockerfile.multistage .

run:
	docker run -d -p 3000:4000 -v vforumdb:/app/data --rm --name cforum iforum

run-dev:
	docker run -d -p 3000:4000 -v vforumdb:/app/data --rm --name cforum iforum

stop:
	docker stop cforum	
	
.DEFAULT_GOAL := build
