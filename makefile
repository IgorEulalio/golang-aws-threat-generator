# Variables
IMAGE = southamerica-east1-docker.pkg.dev/thinking-heaven-421211/golang-event-generator/instrumented
TAG = v2

# Targets
.PHONY: all build docker-build

all: build docker-build

build:
	@echo "Building Go application..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

docker-build: build
	@echo "Building Docker image..."
	docker buildx build --platform linux/amd64 -t $(IMAGE):$(TAG) --push .
