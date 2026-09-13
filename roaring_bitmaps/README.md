# Roaring bitmaps with Go 1.27 vectorisation

Three implementations of the same data structure, sharing one set of SIMD
primitives, so that benchmarking them against each other measures the container
model rather than differences in how well each was optimised.

| Package | What it adds |
|---|---|
| `internal/simdops` | the vector primitives: bitwise ops fused with population count, sorted-array intersection, range edits |
| `v1` | array and bitmap containers |
| `v2` | run containers and `RunOptimize` |
| `v3` | portable-format serialisation, `Rank`/`Select` |
| `bench` | one workload driving all three |

## Building

Everything is behind `//go:build goexperiment.simd && amd64`, with an empty stub
under the negated tag so `go build ./...` at the repository root still works. The
real build needs the Go 1.27 toolchain invoked directly, because the `go` on PATH
is an older launcher that rejects `GOEXPERIMENT` before it switches toolchains:

```
make roaring_test    # go test -race
make roaring_bench   # writes bench_results.txt
make roaring_fuzz
```

The CPU must have AVX-512 F/BW/VL, VPOPCNTDQ and VBMI2; `simdops` panics at init
otherwise. There is no scalar fallback, by design.

## What the vectorisation actually buys

`simd/archsimd` is used rather than the portable `simd` package. The portable API
has no population count and no way to turn a comparison mask into bits, and nearly
every hot Roaring operation needs one or the other, so the portable version would
cover only the pure bitwise loops and drop to `archsimd` everywhere else.

- **Bitmap containers.** `AndTo` and friends load 8 words at a time, apply the
  operation, store it, and accumulate `VPOPCNTQ` of the result in the same pass,
  so a set operation and its cardinality cost one traversal instead of two.
- **Array intersection.** Two strategies, chosen by size ratio: block-wise merge
  with 32 broadcast comparisons and a `VPCOMPRESSW` store per 32-value block, or
  vector-probed galloping when one side is at least 16x larger.
- **Bitmap to array.** `Mask16x32FromBits` plus `Compress` turns 32 bits into
  their positions without a bit-at-a-time loop.
- **Range edits.** Run containers meet bitmaps through `SetRange`, `ClearRange`,
  `FlipRange` and `CopyRange`, which vectorise the whole-word middle of a range.

## What the benchmarks show

From `bench_results.txt` (200k values per dataset):

- Run containers are the whole story on run-heavy and dense data. `And` on the
  dense dataset goes from 17.6us (v1) to 0.44us (v2 after `RunOptimize`), and on
  the runs dataset from 35.5us to 2.8us. `Xor` shows the same shape.
- On sparse data run containers cost slightly more than they save: `RunOptimize`
  finds nothing to collapse, and v2/v3 run a few percent behind v1 through the
  extra dispatch.
- Serialised size is where the encoding shows up most: the dense dataset drops
  from 0.157 to 0.0003 bytes per value after `RunOptimize`, the runs dataset from
  0.53 to 0.008, while sparse and zipf are unchanged at 3.3 and 2.1.
- v3 tracks v2 on set operations, as expected — serialisation does not touch those
  paths. Its own numbers to read are `ToBytes` and `FromBytes`.

## Correctness

Every version is checked against a `map[uint32]struct{}` reference model, in table
tests and in fuzz targets covering set operations, run encoding, serialisation
round-trips and parsing of arbitrary bytes. `bench` additionally asserts that all
three versions produce identical results on every dataset before any timing is
reported.
