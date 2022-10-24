package algo

import "golang.org/x/exp/constraints"

func Less[T constraints.Ordered](a, b T) bool {
	return a < b
}

func Great[T constraints.Ordered](a, b T) bool {
	return a > b
}

func BestNInit[T any](list []T, less Comparator[T], init T) T {
	return Reduce(list, func(ret, v T) T {
		return Tern(less(ret, v), ret, v)
	}, init)
}

func MaxNInit[T any](list []T, c Comparator[T], b Bounds[T]) T {
	return BestNInit(list, c, b.MinValue())
}

func MinNInit[T any](list []T, c Comparator[T], b Bounds[T]) T {
	return BestNInit(list, c, b.MaxValue())
}

func MaxNInitOrdered[T constraints.Ordered](list []T, b Bounds[T]) T {
	return BestNInit(list, Great[T], b.MinValue())
}

func MinNInitOrdered[T constraints.Ordered](list []T, b Bounds[T]) T {
	return BestNInit(list, Less[T], b.MaxValue())
}
