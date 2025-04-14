# Variables
# IMAGE = igoreulalio/golang-threat-generator
IMAGE = us-central1-docker.pkg.dev/crafty-nova-380616/gcf-artifacts/golang-threat-generator
TAG = v2.0.0

# Targets
.PHONY: all build docker-build

all: build docker-build

build:
	@echo "Building Go application..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

docker-build: build
	@echo "Building Docker image..."
	docker buildx build --platform linux/amd64 -t $(IMAGE):$(TAG) --push .

docker-build-gcr: build
	@echo "Building Docker image with GCR docker file..."
	docker buildx build --platform linux/amd64 -f Dockerfile-gcr-agentless -t $(IMAGE):$(TAG) --push .


run-local:
	@echo "Running Go application..."
	go run main.go