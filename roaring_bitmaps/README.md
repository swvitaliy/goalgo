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
make roaring_bench   # 10 runs -> bench_results_count10.txt + bench_results.csv
make roaring_plots   # plots/*.svg via benchdraw
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

`bench/` builds each version the way it is meant to be used — v2 and v3 with
`RunOptimize` applied — and runs all three over the same three datasets: `sparse`
(random values across 2^31, every chunk a tiny array), `runs` (stretches of 512
consecutive values, the shape run containers exist for) and `zipf` (a skewed
distribution that mixes container kinds). Serialisation has no counterpart in
v1/v2, so it is compared against the alternative one would ship without it: a
dump of the sorted `uint32` values, rebuilt on the other side.

From `bench_results.csv` (medians of 10 runs, 200k values per dataset):

- **Run containers are what v2 buys, and only on data that has runs.** On `runs`
  `And` drops from 25.8us (v1) to 1.5us (v2), `AndNot` from 14.4us to 1.7us,
  `Or` from 14.7us to 2.5us. `AndCardinality` does not move (1.36us vs 1.25us):
  v1 already fuses AND and popcount in one SIMD pass without writing a result, so
  there was nothing left to save.
- **They also cost something.** `Contains` on `runs` goes from 5.7ns to 16.9ns —
  a binary search over intervals instead of one bit test. Point lookups pay for
  what set operations gain.
- **On `sparse` v2 and v3 are slower than v1, by ~1.4x on `And`.** Nothing there
  is run-shaped; what shows is the container struct. Carrying a third encoding
  costs one slice header, which puts v2's container in the 64-byte size class
  against v1's 48, and a sparse operation allocates ~32k of them. Packing the
  struct (`card int32`, bitmap as an array pointer) took this gap down from 1.85x;
  the remainder is that size class plus a longer type switch.
- **On `zipf` the three are within noise of each other.**
- **v3 against a raw dump:** on `runs` the stream is 0.008 bytes per value
  against 4, and both directions are two to three orders of magnitude faster
  (`ToBytes` 0.6us vs 218us, `FromBytes` 1.8us vs 1054us). On `sparse` there is
  little to compress (3.3 vs 4 bytes per value) but decoding is still 2.6x faster
  because the format lands directly in containers instead of re-inserting values.

`plot.sh` renders one chart per operation and dataset into `plots/` (the datasets
differ by orders of magnitude, so a shared axis would hide the fast cases). It
needs `benchdraw` built from source: the published module has a broken
dependency, and the fork used here also moves the legend below the plot.

## Correctness

Every version is checked against a `map[uint32]struct{}` reference model, in table
tests and in fuzz targets covering set operations, run encoding, serialisation
round-trips and parsing of arbitrary bytes. `bench` additionally asserts that all
three versions produce identical results on every dataset before any timing is
reported.
