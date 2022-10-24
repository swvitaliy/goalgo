package algo

func Reduce[T any](list []T, acc func(res, v T) T, init T) T {
	res := init
	for _, v := range list {
		res = acc(res, v)
	}
	return res
}

func Tern[T any](cond bool, a, b T) T {
	if cond {
		return a
	} else {
		return b
	}
}
