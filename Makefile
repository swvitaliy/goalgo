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

