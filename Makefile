V := @
OUT_DIR := ./build

.PHONY: vendor
vendor:
	$(V)go mod tidy
	$(V)go mod vendor
	$(V)git add vendor

.PHONY: test
test: GO_TEST_FLAGS += -race -cover
test:
	$(V)go test -mod=vendor $(GO_TEST_FLAGS) --tags=$(GO_TEST_TAGS) ./...