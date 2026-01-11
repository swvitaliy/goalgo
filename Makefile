.PHONY: test
test:
	go test -v -race ./...

.PHONY: fmt
fmt:
	goimports -l -w .

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: bench
bench:
	go test -v -bench=. ./...

.PHONY: fuzz
fuzz:
	go test -v -fuzz=Fuzz -fuzztime=10s ./...:w
