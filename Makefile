# Makefile for ddoser application

APP_NAME := ddoser
SRC_DIR := app
OUT_DIR := bin
COVERPROFILE := covprofile

# Determine the current operating system
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	CGO_ENABLED=0
	OUT_FILE := $(OUT_DIR)/$(APP_NAME)_linux
else
	OUT_FILE := $(OUT_DIR)/$(APP_NAME)
endif

.PHONY: all
all: build

.PHONY: build
build:
	go build -mod=vendor -ldflags "-s -w" -o $(OUT_FILE) ./app

.PHONY: build-linux
build-linux:
	GOOS=linux CGO_ENABLED=0 go build -mod=vendor -ldflags "-s -w" -o $(OUT_FILE)_linux ./app

.PHONY: test
test:
	go test -race -covermode atomic -coverprofile=covprofile ./...

.PHONY: docker-lint
docker-lint:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.53.3 golangci-lint run -v -mod=mod

.PHONY: coverage
coverage: test
	go tool cover -html=$(COVERPROFILE)

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
	git add vendor

.PHONY: clean
clean:
	rm -rf $(OUT_DIR) $(COVERPROFILE)

