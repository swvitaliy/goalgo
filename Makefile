.PHONY: test
test:
	go test -v -race ./...

.PHONY: fmt
fmt:
	goimports -l -w .

.PHONY: lint
lint:
	golangci-lint run ./...
