package typ

import "golang.org/x/exp/constraints"

func BestN[T any, C Comparator[T]](list []T, c C, init T) T {
	return Reduce(list, func(ret, v T) T {
		return Tern(c.Compare(ret, v), ret, v)
	}, init)
}

func MaxN[T any, C Comparator[T]](list []T, c C, b Bounds[T]) T {
	return BestN(list, c, b.MinValue())
}

func MinN[T any, C Comparator[T]](list []T, c C, b Bounds[T]) T {
	return BestN(list, c, b.MaxValue())
}

func MaxNOrdered[T constraints.Ordered](list []T, b Bounds[T]) T {
	return BestN(list, GreatOrdered[T]{}, b.MinValue())
}

func MinNOrdered[T constraints.Ordered](list []T, b Bounds[T]) T {
	return BestN(list, LessOrdered[T]{}, b.MaxValue())
}
