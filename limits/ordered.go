package limits

import "reflect"

const (
	MaxUInt = ^uint(0)
	MaxInt  = int(MaxUInt >> 1)
	MinInt  = -MaxInt - 1
)

func UseOrderedTypes() {
	limits[reflect.Int] = Limits{
		MinValue: MinInt,
		MaxValue: MaxInt,
	}
	// TODO fill ordered type limits
}
