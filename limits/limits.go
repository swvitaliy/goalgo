package limits

import (
	"reflect"
)

type Limits struct {
	MinValue, MaxValue any
}

var limits = make(map[reflect.Type]Limits)

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func MaxValue[T any]() T {
	return limits[typeOf[T]()].MaxValue.(T)
}

func MinValue[T any]() T {
	return limits[typeOf[T]()].MinValue.(T)
}

func AddLimits[T any](mn, mx T) {
	limits[typeOf[T]()] = Limits{MinValue: mn, MaxValue: mx}
}

func init() {
	UseOrderedTypes()
}
