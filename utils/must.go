package utils

import "log"

type Must2Func[T any] func(T, bool) T
type Must2ArgsFunc[T any] func(T, bool) func(...any) T
type Must3Func[T1, T2 any] func(T1, T2, bool) (T1, T2)
type Must3ArgsFunc[T1, T2 any] func(T1, T2, bool) func(...any) (T1, T2)

func MakeMust2Func[T any](pattern string) Must2Func[T] {
	return func(x T, ok bool) T {
		if !ok {
			log.Fatalf(pattern)
		}
		return x
	}
}

func MakeMust2ArgsFunc[T any](pattern string) Must2ArgsFunc[T] {
	return func(x T, ok bool) func(...any) T {
		return func(args ...any) T {
			if !ok {
				log.Fatalf(pattern, args)
			}
			return x
		}
	}
}

func MakeMust3Func[T1, T2 any](pattern string) Must3Func[T1, T2] {
	return func(x T1, y T2, ok bool) (T1, T2) {
		if !ok {
			log.Fatalf(pattern)
		}
		return x, y
	}
}
func MakeMust3ArgsFunc[T1, T2 any](pattern string) Must3ArgsFunc[T1, T2] {
	return func(x T1, y T2, ok bool) func(...any) (T1, T2) {
		return func(args ...any) (T1, T2) {
			if !ok {
				log.Fatalf(pattern, args)
			}
			return x, y
		}
	}
}
