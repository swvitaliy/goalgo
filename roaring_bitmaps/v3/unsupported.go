//go:build !(goexperiment.simd && amd64)

// Package roaringv3 is empty unless built for amd64 with GOEXPERIMENT=simd.
// See the roaring targets in the repository Makefile.
package roaringv3
