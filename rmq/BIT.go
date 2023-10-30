package rmq

import "golang.org/x/exp/constraints"

type BinaryIndexedTree[T any] []T

func BuildBIT[S ~[]T, T any](a S, fn func(a, b T) T) BinaryIndexedTree[T] {
	var b BinaryIndexedTree[T] = make([]T, len(a))
	var zero T
	for i := range b {
		b[i] = zero
	}
	for i := 0; i < len(a); i++ {
		b.Update(i, a[i], fn)
	}
	return b
}

func (b BinaryIndexedTree[T]) Update(i int, v T, fn func(a, b T) T) {
	for ; i < len(b); i = i | (i + 1) {
		b[i] = fn(b[i], v)
	}
}

func (b BinaryIndexedTree[T]) GetRange(l, r int, fn func(a, b T) T, revFn func(a, b T) T) T {
	return revFn(b.Get(r, fn), b.Get(l-1, fn))
}

func (b BinaryIndexedTree[T]) Get(i int, fn func(a, b T) T) T {
	var res T
	for ; i >= 0; i = (i & (i + 1)) - 1 {
		res = fn(res, b[i])
	}
	return res
}

type numeric interface {
	constraints.Integer | constraints.Float
}

type NumericBinaryIndexedTree[T numeric] []T

func BuildSumBIT[T numeric](a []T) NumericBinaryIndexedTree[T] {
	var b NumericBinaryIndexedTree[T] = make([]T, len(a))
	var zero T
	for i := range b {
		b[i] = zero
	}
	for i := 0; i < len(a); i++ {
		b.Update(i, a[i])
	}
	return b
}

func (b NumericBinaryIndexedTree[T]) Update(i int, v T) {
	for ; i < len(b); i = i | (i + 1) {
		b[i] += v
	}
}

func (b NumericBinaryIndexedTree[T]) Sum(i int) T {
	var res T
	for ; i >= 0; i = (i & (i + 1)) - 1 {
		res += b[i]
	}
	return res
}

func (b NumericBinaryIndexedTree[T]) SumRange(l, r int) T {
	return b.Sum(r) - b.Sum(l-1)
}
