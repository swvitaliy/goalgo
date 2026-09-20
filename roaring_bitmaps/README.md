# Roaring bitmaps with Go 1.27 vectorisation

A Roaring bitmap over a small set of SIMD primitives, with a benchmark suite
built to show what each of its optimisations is worth: run containers, the
portable serialisation format, and positional queries.

| Package | What it adds |
|---|---|
| `internal/simdops` | the bit-level primitives: bitwise ops fused with population count, sorted-array intersection, range edits — in a vector build and a scalar one |
| `roaring_bitmaps` (this directory) | the bitmap: array, bitmap and run containers, `RunOptimize`, portable-format serialisation, `Rank`/`Select` |
| `bench` | the workloads, each run with and without `RunOptimize`; a second set of the main operations run on the vector and the scalar `simdops` build |

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

a.SizeInBytes() // heap bytes the bitmap occupies
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
make roaring_bench        # 10 runs -> bench_results_count10.txt + bench_results.csv
make roaring_bench_simd   # scalar + vector passes -> bench_results_simd_count10.txt + bench_results_simd.csv
make roaring_plots        # plots/*.svg via benchdraw
```

`roaring_bench` runs every case ten times at 200ms each on the vector build
(about four minutes) and feeds the raw output through `benchstat -format=csv`,
so the CSV carries medians with 95% confidence intervals. Set `ROARING_N` to change the values per dataset
from the default 200k.

`roaring_bench_simd` runs the `BenchmarkSIMD*` cases twice, first with the
plain `go` (scalar `simdops`) and then with the vector toolchain, into one file.
A binary carries one `simdops` implementation, so this is the only way to get
both on a chart; the cases read the build they run on from `simdops.Vectorized`
and tag themselves `simd=off` or `simd=on`. The main target skips these cases.

Both tools are installed with `go install`:

```
go install golang.org/x/perf/cmd/benchstat@latest
```

`benchdraw` cannot be installed that way: the published module drags in a
dependency at a pseudo-version the proxy no longer serves. Build it from a
checkout instead, dropping its `tools.go` (which only pins the linter it uses on
itself) and bumping the `go` directive so module pruning skips the rest. The
copy used here is also patched to pad the plot area, drop the legend and axis
label, colour each bar separately and label the time axis in ns/us/ms and the
allocation axis in B/KB/MB rather than raw numbers.

`BenchmarkMemory` and `BenchmarkSerializedSize` report `bytes/value` rather
than a time: the bitmap's heap footprint per value (from `Bitmap.SizeInBytes`,
which counts the key and container slices, every container struct and its
payload, and the Rank/Select cache) and the serialised size. `BenchmarkMemory`
measures `zipf` at twice the dataset size. The set operations
carry `B/op` and `allocs/op` from `-benchmem`, and `plots/alloc_*.svg` chart
the bytes each allocates per call.

Case names have the form `runopt=on/data=runs` (`codec=raw/data=runs`,
`method=array/data=runs` for serialisation and Rank/Select, `simd=on/data=runs`
for the vectorisation set), which is what lets both `benchstat -col` and `benchdraw` slice
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

`make roaring_bench_simd` measures this directly: the `BenchmarkSIMD*` cases run
the main operations on the scalar build and then on the vector one, without
`RunOptimize`, so that the `runs` dataset is made of bitmap containers rather
than intervals. From `bench_results_simd.csv` (medians of 10 runs, 200k values
per dataset, `plots/simd_*.svg`):

- **Bitmap containers are where it pays.** On `runs`, `AndCardinality` goes
  from 3.6us to 1.4us (2.5x), `Xor` from 20.5us to 12.1us (1.7x), `And` from
  32.3us to 21.6us (1.5x), `Or`, `AndNot` and `ToArray` gain 1.3x.
- **Array containers gain little.** `zipf` is mostly arrays of a few hundred
  values with a handful of bitmaps at its head: `And` and `AndCardinality`
  gain 1.3x through the vectorised array intersection and the bitmap head,
  `ToArray` 1.2x, while `Or`, `AndNot` and `Xor` come out 1.1–1.3x slower on
  the vector build in this run (the previous run had them even or slightly
  faster; their merges are scalar in both builds and the difference is inside
  run-to-run variation). On `sparse`, 32k arrays of a handful of values each,
  the builds differ by up to 1.15x either way with confidence intervals up to
  35%: the work is walking keys and allocating, not comparing values.
- **The gain is capped by allocation.** `And` on `runs` allocates 39 result
  containers of 8KiB each per call; zeroing and freeing them is the same in
  both builds, which is why the fused loop shows 1.5x where the primitive
  alone would show more. `AndCardinality`, which allocates nothing, is the
  cleanest view of the loop itself.
- **`Contains` is untouched.** A binary search over keys plus one bit test or
  array probe has no loop to vectorise; the two builds are identical.

## What the benchmarks show

`bench/` runs every workload twice over the same three datasets: `runopt=off`
builds the bitmap and leaves it alone, `runopt=on` calls `RunOptimize` first.
`RunOptimize` is the only thing that ever produces run containers, so the
difference between the two is exactly what run containers buy or cost. The
datasets: `sparse` (random values across 2^31, every chunk a tiny array),
`runs` (stretches of 512 consecutive values, the shape run containers exist
for) and `zipf` (stretches of 1 to 64 consecutive values at skewed positions
over a 2^25 key space, the way identifiers are handed out in batches: a
crowded head of bitmap containers, a tail of arrays, and nearly all of it runs
once `RunOptimize` has been called — the one dataset where the three encodings
meet).
Serialisation and Rank/Select have no "off" variant, so each is compared with
the alternative one would keep without it: a sorted `[]uint32` beside the
bitmap (`codec=raw|roaring`, `method=array|roaring`), prepared once rather than
on every call.

From `bench_results.csv` (medians of 10 runs, 200k values per dataset):

- **Run containers are worth an order of magnitude on data that has runs.**
  On `runs`, `And` goes from 20.8us to 1.3us (16x), `AndNot` from 11.7us to
  1.6us (7x), `Or` from 11.9us to 2.0us (6x), `Xor` from 12.0us to 4.2us (3x).
  On `zipf`, whose head is bitmap containers and whose stretches all become
  runs, `And` goes from 209us to 23us (9x), `Or` from 411us to 47us (9x),
  `AndNot` from 279us to 29us (10x), `AndCardinality` from 193us to 22us (9x)
  and `Xor` from 412us to 129us (3x). Operations between two run containers
  stay in interval arithmetic and never materialise the values.
- **Run containers save memory the same way.** `BenchmarkMemory` reports the
  heap footprint from `SizeInBytes`: on `runs` the bitmap shrinks from 0.54 to
  0.013 bytes per value (42x), because a chunk of bitmap words becomes a few
  intervals; on `zipf` (measured at 400k values) from 2.1 to 0.20 bytes per
  value (10x), as bitmaps and arrays alike collapse into a stretch or two per
  chunk. What the set operations allocate per call follows: `And` on `runs`
  from 144KB to 2.5KB, `Or` and `AndNot` from 99KB to 4KB, `Xor` to 9.5KB; on
  `zipf` `And` from 376KB to 46KB, `Or` from 907KB to 88KB, `AndNot` from
  434KB to 85KB, `Xor` from 899KB to 187KB. `sparse` stays at 15 bytes per
  value, most of it the 64-byte container struct and its slot per handful of
  values: no chunk is re-encoded.
- **They also cost something.** `Contains` on `runs` goes from 5.7ns to 15.8ns —
  a binary search over intervals instead of one bit test. Point lookups pay for
  what set operations gain. On `zipf` it goes the other way, 37ns to 22ns: a
  chunk of a few intervals is searched faster than an array of hundreds of
  values.
- **On `sparse` the two variants are the same bitmap.** `RunOptimize` finds
  nothing worth re-encoding. The set operations there allocate 32k containers
  per call, and the allocator dominates the timing: runs of the same bitmap
  differ by up to 1.5x with confidence intervals up to 42%, so nothing on
  `sparse` separates the variants. (Measured earlier against a separate
  two-container implementation, merely carrying the third encoding cost one
  slice header per container — 64 bytes against 48 — visible as 5–35% on
  `sparse`; with a single implementation both variants pay it.)
- **Serialisation wins where run containers do.** On `runs` the stream is
  0.008 bytes per value against 4, writing is 170x faster and reading 100x,
  because containers land in place. On `zipf` it is 0.14 bytes per value
  against 4, writing 10x and reading 4x faster. On `sparse` the sorted copy is
  a plain memory pass and beats the format: writing 4x and reading 8x faster,
  with 0.7 bytes per value saved against 4.
- **Rank/Select trade speed for memory.** The sorted copy answers `Select` in
  ~1.5ns against 13–68ns — an index is unbeatable — but costs 4 bytes per
  value where the prefix-sum cache costs one int per container (0.0006–1.3
  bytes per value). `Rank` favours the bitmap wherever containers are few:
  2.4x faster than the copy on `runs`, 1.5x on `zipf`, 1.3x slower on
  `sparse`.

## Correctness

The bitmap is checked against a `map[uint32]struct{}` reference model, in table
tests and in fuzz targets covering set operations, run encoding, serialisation
round-trips, Rank/Select and parsing of arbitrary bytes. The same suite runs on
the vector and the scalar `simdops` build. `bench` additionally asserts, before
any timing is reported, that the bitmap answers the same with and without
`RunOptimize`, that both codecs round-trip to the same values, and that the
array baseline and the bitmap agree on every Rank and Select it measures. For
the vectorisation set it checks, on whichever build is compiled in, that the
operations it times satisfy the set identities that tie them together
(`|A|+|B| = |A∪B|+|A∩B|`, `|A∖B| = |A|-|A∩B|`, `|A⊕B| = |A∪B|-|A∩B|`,
`AndCardinality` equal to the cardinality of `And`) and that `ToArray` lists
exactly the set, so the two passes are known to time the same answers.
