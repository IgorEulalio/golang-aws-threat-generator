# Variables
IMAGE = igoreulalio/release-book-architecture
TAG = v1.0.9

# Targets
.PHONY: all build docker-build

all: build docker-build

build:
	@echo "Building Go application..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

docker-build: build
	@echo "Building Docker image..."
	docker buildx build --platform linux/amd64 -t $(IMAGE):$(TAG) --push .
