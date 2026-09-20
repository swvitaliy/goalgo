#!/usr/bin/env bash
# Renders comparison plots from a benchmark run with benchdraw.
#
# Usage: roaring_bitmaps/plot.sh [input.txt] [outdir] [simd_input.txt]
#
# Benchmark cases are named runopt=<on|off>/data=<shape> (codec=raw|roaring and
# method=array|roaring for serialisation and Rank/Select), which is what lets benchdraw
# slice one run by dimension. The SIMD benchmarks come from a separate file with
# cases named simd=<off|on>/data=<shape>, one pass per simdops build.
#
# There is one chart per operation and dataset rather than one per operation with
# the datasets side by side: timings differ by orders of magnitude between the
# datasets, and on a shared linear axis the fast cases collapse into the baseline
# and show nothing.
set -euo pipefail

input="${1:-roaring_bitmaps/bench_results_count10.txt}"
outdir="${2:-roaring_bitmaps/plots}"
simd_input="${3:-roaring_bitmaps/bench_results_simd_count10.txt}"
mkdir -p "$outdir"
rm -f "$outdir"/*.svg

draw() {
	local name="$1" filter="$2" x="$3" y="$4" title="$5" from="${6:-$input}"
	benchdraw -filter="$filter" -x="$x" -y="$y" \
		-title="$title" -input="$from" -output="$outdir/$name.svg"
	echo "  $outdir/$name.svg"
}

echo "plots from $input:"

for op in Build Contains And Or AndNot Xor AndCardinality; do
	for data in sparse runs zipf; do
		draw "${op,,}_${data}" "Benchmark$op/data=$data" runopt ns/op "$op, $data"
	done
done

for data in sparse runs zipf; do
	draw "size_${data}" "BenchmarkSerializedSize/data=$data" codec bytes/value "Serialised size, $data"
	draw "tobytes_${data}" "BenchmarkToBytes/data=$data" codec ns/op "ToBytes, $data"
	draw "frombytes_${data}" "BenchmarkFromBytes/data=$data" codec ns/op "FromBytes, $data"
done

for data in sparse runs zipf; do
	draw "rank_${data}" "BenchmarkRank/data=$data" method ns/op "Rank, $data"
	draw "select_${data}" "BenchmarkSelect/data=$data" method ns/op "Select, $data"
	draw "aux_${data}" "BenchmarkRank/data=$data" method aux-bytes/value "Auxiliary state, $data"
done

# What vectorisation buys: the same operation and data, primitives on the x axis.
if [[ -f "$simd_input" ]]; then
	echo "plots from $simd_input:"
	for op in Contains And Or AndNot Xor AndCardinality ToArray; do
		for data in sparse runs zipf; do
			draw "simd_${op,,}_${data}" "BenchmarkSIMD$op/data=$data" simd ns/op "$op, $data: scalar vs AVX-512" "$simd_input"
		done
	done
fi
