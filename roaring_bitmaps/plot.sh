#!/usr/bin/env bash
# Renders comparison plots from a benchmark run with benchdraw.
#
# Usage: roaring_bitmaps/plot.sh [input.txt] [outdir]
#
# Benchmark cases are named version=<v>/data=<shape> (codec=<c>/data=<shape> for
# serialisation), which is what lets benchdraw slice one run by dimension.
#
# There is one chart per operation and dataset rather than one per operation with
# the datasets side by side: timings differ by orders of magnitude between the
# datasets, and on a shared linear axis the fast cases collapse into the baseline
# and show nothing.
set -euo pipefail

input="${1:-roaring_bitmaps/bench_results_count10.txt}"
outdir="${2:-roaring_bitmaps/plots}"
mkdir -p "$outdir"
rm -f "$outdir"/*.svg

draw() {
	local name="$1" filter="$2" x="$3" y="$4" title="$5"
	benchdraw -filter="$filter" -x="$x" -y="$y" \
		-title="$title" -input="$input" -output="$outdir/$name.svg"
	echo "  $outdir/$name.svg"
}

echo "plots from $input:"

# Version comparison, on the vectorised build.
for op in Build Contains And Or AndNot Xor AndCardinality; do
	for data in sparse runs zipf; do
		draw "${op,,}_${data}" "Benchmark$op/data=$data/simd=on" version ns/op "$op, $data"
	done
done

for data in sparse runs zipf; do
	draw "size_${data}" "BenchmarkSerializedSize/data=$data/simd=on" codec bytes/value "Serialised size, $data"
	draw "tobytes_${data}" "BenchmarkToBytes/data=$data/simd=on" codec ns/op "ToBytes, $data"
	draw "frombytes_${data}" "BenchmarkFromBytes/data=$data/simd=on" codec ns/op "FromBytes, $data"
done

# What vectorisation buys: the same version and data, primitives on the x axis.
# v1 is the interesting one — its bitmap containers are where the SIMD loops run;
# v2 on runs is included to show that interval arithmetic gains nothing.
for op in And Or AndNot Xor Contains; do
	for data in sparse runs zipf; do
		draw "simd_${op,,}_${data}" "Benchmark$op/version=v1/data=$data" simd ns/op "$op, $data, v1: scalar vs AVX-512"
	done
done
for op in And Or; do
	draw "simd_${op,,}_runs_v2" "Benchmark$op/version=v2/data=runs" simd ns/op "$op, runs, v2: scalar vs AVX-512"
done
