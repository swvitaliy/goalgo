package goalgo

import (
	"math/bits"
)

type SparseTable [][]int

func __lg(x uint64) int {
	// https://stackoverflow.com/questions/11376288/fast-computing-of-log2-for-64-bit-integers
	return 8*64 - bits.LeadingZeros64(x) - 1
}

func NewSparseTable(n int) SparseTable {
	logN := __lg(uint64(n))
	a := make([][]int, logN)
	for i := range a {
		a[i] = make([]int, n)
	}
	return a
}

func BuildSparseTable(a []int) SparseTable {
	st := NewSparseTable(len(a))
	n := len(st[0])
	logN := __lg(uint64(n))
	for i := 0; i < n; i++ {
		st[0][i] = a[i]
	}
	for k := 0; k < logN; k++ {
		for i := 0; i < n-1<<k; i++ {
			st[k+1][i] = min(st[k][i], st[k][i+1])
		}
	}
	return st
}

func (stp *SparseTable) Rmq(l, r int) int {
	st := *stp
	k := __lg(uint64(r - l))
	return min(st[k][l], st[k][r-(1<<k)])
}
