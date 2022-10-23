package typ

import (
	"math"
)

type Bounds[T any] interface {
	MinValue() T
	MaxValue() T
}

type IntBoundsImpl struct {
}

func (b IntBoundsImpl) MaxValue() int {
	return math.MaxInt
}

func (b IntBoundsImpl) MinValue() int {
	return math.MinInt
}

var IntBounds Bounds[int] = IntBoundsImpl{}

// TODO Add more bounds
