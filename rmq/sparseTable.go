package rmq

import "math/bits"

type SparseTable[T any] [][]T

func lg(x uint64) int {
	return 63 - bits.LeadingZeros64(x)
}

func NewSparseTable[T any](n int) SparseTable[T] {
	logN := lg(uint64(n)) + 1
	a := make([][]T, logN)
	for i := range a {
		a[i] = make([]T, n)
	}
	return a
}

func BuildSparseTable[S ~[]T, T any](a S, fn func(a, b T) T) SparseTable[T] {
	st := NewSparseTable[T](len(a))
	if len(st) == 0 {
		return st
	}
	n := len(st[0])
	logN := lg(uint64(n)) + 1
	copy(st[0], a)
	for k := 0; k < logN-1; k++ {
		for i := 0; i <= n-(2<<k); i++ {
			st[k+1][i] = fn(
				st[k][i],
				st[k][i+(1<<k)])
		}
	}
	return st
}

func (stp *SparseTable[T]) Rmq(l, r int, fn func(a, b T) T) T {
	st := *stp
	k := lg(uint64(r - l))
	return fn(
		st[k][l],
		st[k][r-(1<<k)+1])
}
