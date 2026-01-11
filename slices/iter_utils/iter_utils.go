package iter_utils

import (
	"iter"
	"time"

	"golang.org/x/exp/rand"
)

func RandomInt64Iter(K int, min, max int64) iter.Seq[int64] {
	if max <= min {
		panic("max must be greater than min")
	}

	return func(yield func(int64) bool) {
		rand.Seed(uint64(time.Now().UnixNano()))

		for i := 0; i < K; i++ {
			v := rand.Int63n(max-min+1) + min
			if !yield(v) {
				return
			}
		}
	}
}

func ToSlice[T any](it iter.Seq[T]) []T {
	s := make([]T, 0)
	for v := range it {
		s = append(s, v)
	}
	return s
}

func PartiallyShuffledSlice(arr []int64, shuffleFraction float64) []int64 {
	rand.Seed(uint64(time.Now().UnixNano()))
	n := len(arr)
	numToShuffle := int(float64(n) * shuffleFraction)
	for i := 0; i < numToShuffle; i++ {
		a := rand.Intn(n)
		b := rand.Intn(n)
		arr[a], arr[b] = arr[b], arr[a]
	}

	return arr
}
