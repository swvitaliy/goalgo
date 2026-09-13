//go:build !(goexperiment.simd && amd64)

// Package simdops is empty unless built for amd64 with GOEXPERIMENT=simd.
//
// The Roaring bitmaps in this directory are a study of Go 1.27 vectorisation and
// carry no scalar fallback by design, so the package compiles to nothing rather
// than breaking `go build ./...` at the repository root. Use the roaring targets
// in the Makefile to build and test it for real.
package simdops
