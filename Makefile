GO=GO111MODULE=on go

.PHONY: test
test:
	$(GO) test ./pkg/...

.PHONY: lint
lint:
	golangci-lint run --fix --disable-all -E goimports

