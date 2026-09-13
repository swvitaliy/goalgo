# GoAlgo

A bunch of algorithms in golang.

| Package | |
|---|---|
| `bloom` | Bloom filter |
| `cms` | count-min sketch |
| `graphs` | cycle detection, shortest paths, minimum spanning tree, strongly connected components |
| `hashes` | consistent hashing and rendezvous hashing, with benchmarks and plots |
| `heap` | heaps, including a pairing heap |
| `intervals` | insert, intersect and merge over interval lists |
| `lca` | lowest common ancestor, disjoint-set union, Tarjan |
| `lfu`, `lru` | cache eviction policies |
| `limits` | min/max values per type, for generic code |
| `linked_list`, `skip_list` | lists |
| `rmq` | range minimum query: Fenwick tree, segment tree, sparse table |
| `roaring_bitmaps` | Roaring bitmap on Go 1.27 SIMD, with a benchmark report — see its [README](roaring_bitmaps/README.md) |
| `subseq` | longest common subsequence, longest increasing subsequence |
| `trees` | binary search tree |
| `slices`, `sorts`, `strings` | the usual helpers |

Most packages build with a plain `go test ./...`. `roaring_bitmaps` also has a
vector build behind `GOEXPERIMENT=simd`; see `make roaring_test` and friends in
the Makefile.
