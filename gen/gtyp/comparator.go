package gtyp

import "golang.org/x/exp/constraints"

type FuncComparator[T any] func(x, y T) bool

type Comparator[T any] interface {
	~struct{}
	Compare(x, y T) bool
}

type GreatOrdered[T constraints.Ordered] struct{}

func (c GreatOrdered[T]) Compare(x, y T) bool {
	return x > y
}

type LessOrdered[T constraints.Ordered] struct{}

func (c LessOrdered[T]) Compare(x, y T) bool {
	return x < y
}

type EqualOrdered[T constraints.Ordered] struct{}

func (c EqualOrdered[T]) Compare(x, y T) bool {
	return x == y
}
