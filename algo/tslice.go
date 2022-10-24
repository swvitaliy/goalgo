package algo

import "sort"

type TSlice[T any] struct {
	a []T
	c Comparator[T]
}

func (x TSlice[T]) Len() int           { return len(x.a) }
func (x TSlice[T]) Less(i, j int) bool { return x.c(x.a[i], x.a[j]) }
func (x TSlice[T]) Swap(i, j int)      { x.a[i], x.a[j] = x.a[j], x.a[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x TSlice[T]) Sort() { sort.Sort(x) }

func (x TSlice[T]) Data() []T { return x.a }
