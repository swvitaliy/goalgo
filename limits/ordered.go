package limits

import "reflect"

const (
	MaxUInt = ^uint(0)
	MaxInt  = int(MaxUInt >> 1)
	MinInt  = -MaxInt - 1
)

func UseOrderedTypes() {
	AddKindLimits(reflect.Int, Limits{
		MinValue: MinInt,
		MaxValue: MaxInt,
	})
	AddKindLimits(reflect.Uint, Limits{
		MinValue: 0,
		MaxValue: MaxUInt,
	})

	// TODO fill ordered type limits
}
