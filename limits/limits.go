package limits

import (
	"reflect"
)

type Limits struct {
	MinValue, MaxValue any
}

var limits = make(map[reflect.Kind]Limits)

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func MaxValue[T any]() T {
	return limits[typeOf[T]().Kind()].MaxValue.(T)
}

func MinValue[T any]() T {
	return limits[typeOf[T]().Kind()].MinValue.(T)
}

func AddKindLimits(kind reflect.Kind, l Limits) {
	limits[kind] = l
}

func AddTypeLimits[T any](l Limits) {
	limits[typeOf[T]().Kind()] = l
}

func init() {
	UseOrderedTypes()
}
