package misc

import "cmp"

type PriorPair[P cmp.Ordered, T comparable] struct {
	Prior P
	Val   T
}

// PriorPairCmp return -1 if a < b, 1 if a > b, 0 if a == b
// When a and b have same priority but different values return 1
func PriorPairCmp[P cmp.Ordered, T comparable](a, b PriorPair[P, T]) int {
	if a.Prior < b.Prior {
		return -1
	}
	if a.Prior > b.Prior {
		return 1
	}
	if a.Val != b.Val {
		return 1
	}
	return 0 // priority and value are equal
}

func MakePriorPair[P cmp.Ordered, T comparable](p P, v T) PriorPair[P, T] {
	return PriorPair[P, T]{p, v}
}
