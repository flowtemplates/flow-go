.DEFAULT_GOAL := test

.PHONY: fmt
fmt:
	@gofumpt -w -l .
	@goimports -w -l .

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@gotestsum -f testdox
