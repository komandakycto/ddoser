V := @
OUT_DIR := ./build
MAIN_PKG := github.com/komandakycto/ddoser

.PHONY: vendor
vendor:
	$(V)go mod tidy
	$(V)go mod vendor
	$(V)git add vendor

.PHONY: test
test: GO_TEST_FLAGS += -race -cover
test:
	$(V)go test -mod=vendor $(GO_TEST_FLAGS) --tags=$(GO_TEST_TAGS) ./...

.PHONY: build
build:
	@echo BUILDING $(OUT_DIR)/ddoser
	$(V)go build -mod=vendor -ldflags "-s -w" -o $(OUT_DIR)/ddoser $(MAIN_PKG)/app/main.go
	@echo DONE

.PHONY: linux
linux: export GOOS := linux
linux: export GOARCH := amd64
linux: export CGO_ENABLED := 0
linux: build