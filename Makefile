# Variables
IMAGE_NAME := reddit-app
TAG := latest

# Targets
build:
    go build -o reddit-app .

test:
    go test ./...

docker-build:
    docker build -t $(IMAGE_NAME):$(TAG) .

docker-run:
    docker run -p 8080:8080 $(IMAGE_NAME):$(TAG)

clean:
    rm -f reddit-app

.PHONY: build test docker-build docker-run clean
