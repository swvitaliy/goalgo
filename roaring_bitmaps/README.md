# Roaring bitmaps with Go 1.27 vectorisation

A Roaring bitmap over a small set of SIMD primitives, with a benchmark suite
built to show what each of its optimisations is worth: run containers, the
portable serialisation format, and positional queries.

| Package | What it adds |
|---|---|
| `internal/simdops` | the bit-level primitives: bitwise ops fused with population count, sorted-array intersection, range edits — in a vector build and a scalar one |
| `roaring_bitmaps` (this directory) | the bitmap: array, bitmap and run containers, `RunOptimize`, portable-format serialisation, `Rank`/`Select` |
| `bench` | the workloads, each run with and without `RunOptimize` |

## Usage

```go
import roaring "goalgo/roaring_bitmaps"

a := roaring.New(1, 2, 3, 100_000, 100_001)
b := roaring.New(2, 3, 4, 100_001, 1<<20)

union := a.Or(b)  // {1 2 3 4 100000 100001 1048576}
inter := a.And(b) // {2 3 100001}
a.AndNot(b)       // values of a not in b
a.Xor(b)          // symmetric difference

a.AndCardinality(b) // 3, without building the intersection
a.Intersects(b)     // true, stops at the first shared value
a.Contains(100_001) // true

// Run encoding is never chosen on its own. Once a bitmap holds runs of
// consecutive values, ask for it explicitly; set operations between two
// run containers then stay in interval arithmetic.
a.RunOptimize()

data := a.ToBytes()              // portable Roaring format
c, err := roaring.FromBytes(data) // c.ToArray() == a.ToArray()

a.Rank(100_000) // 4: how many values are <= 100000
a.Select(3)     // 100000, true: the 4th smallest value
```

`Or`, `And`, `AndNot` and `Xor` return a new bitmap and leave both operands
untouched. The zero value of `Bitmap` is an empty set ready for use; a `Bitmap`
is not safe for concurrent modification.

## How a set operation finds its work

A value's high 16 bits pick a chunk, its low 16 bits a position inside that
chunk's container. A bitmap keeps its chunk keys in a sorted `[]uint16` with the
containers in a parallel slice, so a binary operation never searches: it walks
both key lists with two indices, like a merge join, pairing equal keys and
handling one-sided ones according to the operation (`Or` clones them, `And`
skips them). Empty results drop out together with their key.

Each pair of containers is dispatched on its two encodings. Array against array
goes through `simdops.IntersectArrays` or a merge; bitmap against bitmap through
the fused `AndTo`/`OrTo` loops; run against run stays in interval arithmetic; the
mixed pairs walk the cheaper structure — array values probed against runs or
bits, bitmaps edited range by range. The result is then normalised to the
cheapest encoding for its cardinality: an array up to 4096 values, a bitmap
above that, runs only if it was already runs and still the smallest.

Point operations — `Add`, `Remove`, `Contains` — find their chunk with one
binary search over the keys.

The bitmap never calls `simdops` directly. Every primitive it needs is listed in
the `Ops` interface (`ports.go`), and each `Bitmap` carries the implementation
it runs on: `New` and `FromBytes` use `simdops.Ops{}`, `NewWithOps` and
`FromBytesWithOps` take any other, and every bitmap derived by `And`, `Or`,
`Clone` and friends inherits it. Containers do not store it — that would grow
the 64-byte struct — so the container functions take it as their first
argument. `ports_mock.go` is the gomock double, regenerated with `go generate`.

## Building

`simdops` has two implementations of its hot primitives, chosen by build tag:
`vector_amd64.go` under `GOEXPERIMENT=simd` on amd64, and `scalar.go` everywhere
else. Both export the same API and pass the same tests; `simdops.Vectorized`
says which one is compiled in. The vector build also exports
`simdops.VectorOps` (`vector_ops.go`), a struct whose methods forward to the
AVX-512 primitives. Everything above `simdops` is tag-free, so a
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
copy used here is also patched to pad the plot area, drop the legend and axis
label, colour each bar separately and label the time axis in ns/us/ms rather
than raw nanoseconds.

Case names have the form `runopt=on/data=runs` (`codec=raw/data=runs`,
`method=array/data=runs` for serialisation and Rank/Select), which is what lets both `benchstat -col` and `benchdraw` slice
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

`bench/` runs every workload twice over the same three datasets: `runopt=off`
builds the bitmap and leaves it alone, `runopt=on` calls `RunOptimize` first.
`RunOptimize` is the only thing that ever produces run containers, so the
difference between the two is exactly what run containers buy or cost. The
datasets: `sparse` (random values across 2^31, every chunk a tiny array),
`runs` (stretches of 512 consecutive values, the shape run containers exist
for) and `zipf` (a skewed distribution that mixes container kinds).
Serialisation and Rank/Select have no "off" variant, so each is compared with
the alternative one would keep without it: a sorted `[]uint32` beside the
bitmap (`codec=raw|roaring`, `method=array|roaring`), prepared once rather than
on every call.

From `bench_results.csv` (medians of 10 runs, 200k values per dataset):

- **Run containers are worth an order of magnitude on data that has runs.**
  On `runs`, `And` goes from 23.9us to 1.5us (16x), `AndNot` from 14.5us to
  1.9us (8x), `Or` from 14.3us to 2.4us (6x), `Xor` from 14.8us to 5.4us (3x).
  Operations between two run containers stay in interval arithmetic and never
  materialise the values.
- **They also cost something.** `Contains` on `runs` goes from 5.8ns to 16.5ns —
  a binary search over intervals instead of one bit test. Point lookups pay for
  what set operations gain.
- **On `sparse` and `zipf` the two variants are the same bitmap.** `RunOptimize`
  finds nothing worth re-encoding, and every difference (0.8–1.2x) is inside the
  15–23% confidence intervals. (Measured earlier against a separate
  two-container implementation, merely carrying the third encoding cost one
  slice header per container — 64 bytes against 48 — visible as 5–35% on
  `sparse`; with a single implementation both variants pay it.)
- **Serialisation wins only where run containers do.** On `runs` the stream is
  0.008 bytes per value against 4, writing is 180x faster and reading 120x,
  because containers land in place. On `sparse` and `zipf` the sorted copy is a
  plain memory pass and beats the format: writing 3.6x and reading 6x faster on
  `sparse`, 1.2x and 1.8x on `zipf`, with 3.3 and 2.1 bytes per value saved
  against 4.
- **Rank/Select trade speed for memory.** The sorted copy answers `Select` in
  ~1.5ns against 14–76ns — an index is unbeatable — but costs 4 bytes per value
  where the prefix-sum cache costs one int per container (0.0006–1.3 bytes per
  value). `Rank` is even: 1.1–1.3x slower than the copy on `zipf` and `sparse`,
  twice as fast on `runs`.

## Correctness

The bitmap is checked against a `map[uint32]struct{}` reference model, in table
tests and in fuzz targets covering set operations, run encoding, serialisation
round-trips, Rank/Select and parsing of arbitrary bytes. The same suite runs on
the vector and the scalar `simdops` build. `bench` additionally asserts, before
any timing is reported, that the bitmap answers the same with and without
`RunOptimize`, that both codecs round-trip to the same values, and that the
array baseline and the bitmap agree on every Rank and Select it measures.
