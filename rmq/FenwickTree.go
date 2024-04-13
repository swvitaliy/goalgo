package rmq

import (
	"goalgo/slices"
	"golang.org/x/exp/constraints"
)

// Binary Indexed Tree

type FenwickTree[T any] []T

func BuildFenwickTree[S ~[]T, T any](a S, fn func(a, b T) T) FenwickTree[T] {
	var b FenwickTree[T] = make([]T, len(a))
	var zero T
	slices.Fill(b, zero)
	for i := 0; i < len(a); i++ {
		b.Update(i, a[i], fn)
	}
	return b
}

func (b FenwickTree[T]) Update(i int, v T, fn func(a, b T) T) {
	for ; i < len(b); i = i | (i + 1) {
		b[i] = fn(b[i], v)
	}
}

func (b FenwickTree[T]) GetRange(l, r int, fn func(a, b T) T, revFn func(a, b T) T) T {
	return revFn(b.Get(r, fn), b.Get(l-1, fn))
}

func (b FenwickTree[T]) Get(i int, fn func(a, b T) T) T {
	var res T
	for ; i >= 0; i = (i & (i + 1)) - 1 {
		res = fn(res, b[i])
	}
	return res
}

type numeric interface {
	constraints.Integer | constraints.Float
}

type NumericFenwickTree[T numeric] []T

func BuildSumFenwickTree[T numeric](a []T) NumericFenwickTree[T] {
	var b NumericFenwickTree[T] = make([]T, len(a))
	var zero T
	slices.Fill(b, zero)
	for i := 0; i < len(a); i++ {
		b.Update(i, a[i])
	}
	return b
}

func (b NumericFenwickTree[T]) Update(i int, v T) {
	for ; i < len(b); i = i | (i + 1) {
		b[i] += v
	}
}

func (b NumericFenwickTree[T]) Sum(i int) T {
	var res T
	for ; i >= 0; i = (i & (i + 1)) - 1 {
		res += b[i]
	}
	return res
}

func (b NumericFenwickTree[T]) SumRange(l, r int) T {
	return b.Sum(r) - b.Sum(l-1)
}
