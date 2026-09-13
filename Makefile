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

bench_hashes_100k:
	N=100000 go test -bench=. -benchmem ./hashes | tee ./hashes/bench_results_100k.txt

hashes_bench_1M:
	N=1000000 go test -bench=. -benchmem ./hashes | tee ./hashes/bench_results_1M.txt

hashes_bench_10M:
	N=10000000 go test -bench=. -benchmem ./hashes | tee ./hashes/bench_results_10M.txt

hashes_plot:
	cd ./hashes && go run ./plot/plot.go

skiplist_bench:
	go test -bench=. -benchmem ./skiplist | tee ./skiplist/bench_results.txt


# The go binary on PATH is an older launcher that rejects GOEXPERIMENT before it
# switches toolchains, so the resolved 1.27 toolchain is invoked directly.
GO_SIMD_ROOT := $(shell go env GOROOT)
GO_SIMD := GOROOT=$(GO_SIMD_ROOT) GOTOOLCHAIN=local GOEXPERIMENT=simd $(GO_SIMD_ROOT)/bin/go

.PHONY: roaring_test
roaring_test:
	$(GO_SIMD) test -race ./roaring_bitmaps/...

.PHONY: roaring_bench
roaring_bench:
	$(GO_SIMD) test -bench=. -benchmem ./roaring_bitmaps/... | tee ./roaring_bitmaps/bench_results.txt

.PHONY: roaring_fuzz
roaring_fuzz:
	$(GO_SIMD) test -run=Fuzz -fuzz=Fuzz -fuzztime=30s ./roaring_bitmaps/...
