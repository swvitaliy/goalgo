package gtyp

import "sort"

type TSlice[T any, C Comparator[T]] []T

func (x TSlice[T, C]) Len() int           { return len(x) }
func (x TSlice[T, C]) Less(i, j int) bool { return C{}.Compare(x[i], x[j]) } // Need to benchmark it!
func (x TSlice[T, C]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x TSlice[T, C]) Sort() { sort.Sort(x) }
