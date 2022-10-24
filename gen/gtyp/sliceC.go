package gtyp

import "sort"

type TSliceC[T any] struct {
	a []T
	c FuncComparator[T]
}

func (x TSliceC[T]) Len() int           { return len(x.a) }
func (x TSliceC[T]) Less(i, j int) bool { return x.c(x.a[i], x.a[j]) }
func (x TSliceC[T]) Swap(i, j int)      { x.a[i], x.a[j] = x.a[j], x.a[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x TSliceC[T]) Sort() { sort.Sort(x) }
