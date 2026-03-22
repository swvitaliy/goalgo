.PHONY: test
test:
	go test -v -race ./...

.PHONY: fmt
fmt:
	goimports -l -w .

.PHONY: lint
lint:
	golangci-lint run ./...

bench_hashes:
	go test -bench=. -benchmem ./hashes | tee ./hashes/bench_results.txt
