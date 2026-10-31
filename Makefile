GO ?= go
VERSION ?= dev
IMAGE ?= go-cli-template:local

.PHONY: help run test build fmt vet check clean docker-build

help:
	@printf '%s\n' 'run          Run CLI (ARGS="greet --name Go")' 'test         Run tests with race detection' 'build        Build bin/go-cli-template' 'fmt          Format Go source' 'vet          Run go vet' 'check        Check formatting, vet and tests' 'clean        Remove build output' 'docker-build Build container image'

run:
	$(GO) run ./cmd/go-cli-template $(ARGS)

test:
	$(GO) test -race -cover ./...

build:
	$(GO) build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o bin/go-cli-template ./cmd/go-cli-template

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

check: vet test
	@test -z "$$($(GO) fmt ./...)" || (echo 'Go files needed formatting; run make fmt before committing'; exit 1)

clean:
	rm -rf bin

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE) -f docker/Dockerfile .
