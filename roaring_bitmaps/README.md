# Roaring bitmaps with Go 1.27 vectorisation

Three implementations of the same data structure, sharing one set of SIMD
primitives, so that benchmarking them against each other measures the container
model rather than differences in how well each was optimised.

| Package | What it adds |
|---|---|
| `internal/simdops` | the bit-level primitives: bitwise ops fused with population count, sorted-array intersection, range edits — in a vector build and a scalar one |
| `v1` | array and bitmap containers |
| `v2` | run containers and `RunOptimize` |
| `v3` | portable-format serialisation, `Rank`/`Select` |
| `bench` | one workload driving all three |

## Building

`simdops` has two implementations of its hot primitives, chosen by build tag:
`vector_amd64.go` under `GOEXPERIMENT=simd` on amd64, and `scalar.go` everywhere
else. Both export the same API and pass the same tests; `simdops.Vectorized`
says which one is compiled in. Everything above `simdops` is tag-free, so a
plain `go test ./...` at the repository root exercises the scalar build.

The vector build needs the Go 1.27 toolchain invoked directly, because the `go`
on PATH is an older launcher that rejects `GOEXPERIMENT` before it switches
toolchains:

```
make roaring_test    # go test -race
make roaring_fuzz
```

For the vector build the CPU must have AVX-512 F/BW/VL, VPOPCNTDQ and VBMI2;
`simdops` panics at init otherwise. The scalar build has no such requirement.

## Reproducing the benchmarks

```
make roaring_bench   # 10 runs -> bench_results_count10.txt + bench_results.csv
make roaring_plots   # plots/*.svg via benchdraw
```

`roaring_bench` runs every case ten times at 200ms each on the vector build
(about four minutes) and feeds the raw output through `benchstat -format=csv`,
so the CSV carries medians with 95% confidence intervals. Set `ROARING_N` to change the values per dataset
from the default 200k.

Both tools are installed with `go install`:

```
go install golang.org/x/perf/cmd/benchstat@latest
```

`benchdraw` cannot be installed that way: the published module drags in a
dependency at a pseudo-version the proxy no longer serves. Build it from a
checkout instead, dropping its `tools.go` (which only pins the linter it uses on
itself) and bumping the `go` directive so module pruning skips the rest. The
copy used here is also patched to pad the plot area, drop the legend and label
the time axis in ns/us/ms rather than raw nanoseconds.

Case names have the form `version=v2/data=runs` (`codec=raw/data=runs` for
serialisation), which is what lets both `benchstat -col` and `benchdraw` slice
a single run by dimension.

## What the vectorisation actually buys

In the vector build, `simd/archsimd` is used rather than the portable `simd` package. The portable API
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
  `And` drops from 26.8us (v1) to 1.5us (v2), `AndNot` from 14.5us to 1.8us,
  `Or` from 13.8us to 2.2us. `AndCardinality` does not move (1.51us vs 1.33us):
  v1 already fuses AND and popcount in one SIMD pass without writing a result, so
  there was nothing left to save.
- **They also cost something.** `Contains` on `runs` goes from 5.6ns to 17.4ns —
  a binary search over intervals instead of one bit test. Point lookups pay for
  what set operations gain.
- **On `sparse` v2 and v3 are slower than v1, by up to 1.5x on `And` depending on
  the run (the confidence intervals there reach 30%).** Nothing there
  is run-shaped; what shows is the container struct. Carrying a third encoding
  costs one slice header, which puts v2's container in the 64-byte size class
  against v1's 48, and a sparse operation allocates ~32k of them. Packing the
  struct (`card int32`, bitmap as an array pointer) took this gap down from 1.85x;
  the remainder is that size class plus a longer type switch.
- **On `zipf` the three are within noise of each other.**
- **v3 against a raw dump:** on `runs` the stream is 0.008 bytes per value
  against 4, and both directions are two to three orders of magnitude faster
  (`ToBytes` 0.6us vs 226us, `FromBytes` 1.9us vs 1152us). On `sparse` there is
  little to compress (3.3 vs 4 bytes per value) but decoding is still 2.7x faster
  because the format lands directly in containers instead of re-inserting values.

`report.html` is a self-contained page — findings plus every chart inlined —
that opens straight from disk. `plot.sh` renders one chart per operation and dataset into `plots/` (the datasets
differ by orders of magnitude, so a shared axis would hide the fast cases).

## Correctness

Every version is checked against a `map[uint32]struct{}` reference model, in table
tests and in fuzz targets covering set operations, run encoding, serialisation
round-trips and parsing of arbitrary bytes. `bench` additionally asserts that all
three versions produce identical results on every dataset before any timing is
reported.
