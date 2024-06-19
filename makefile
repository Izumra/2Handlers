SHELL=/usr/bin/bash

build:
	docker build -t twohandlers:multistage -f Dockerfile.multistage ./

run:
	docker compose up

docs:
	swag init -g ./cmd/main.go
