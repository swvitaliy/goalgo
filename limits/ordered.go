package limits

const (
	MaxUInt = ^uint(0)
	MaxInt  = int(MaxUInt >> 1)
	MinInt  = -MaxInt - 1
)

func UseOrderedTypes() {
	AddLimits[int](MinInt, MaxInt)
	AddLimits[uint](0, MaxUInt)

	// TODO fill ordered type limits
}
