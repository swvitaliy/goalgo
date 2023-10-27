package sorts

import (
	"goalgo/slice"
	"testing"
)

func TestBitonicSort(t *testing.T) {
	a := []int{1, 4, 6, 6, 32, 2, 7, 3, 34, 6, 6}
	BitonicSort(a, true)
	if !slice.IsSorted(a, true) {
		t.Fail()
	}
	BitonicSort(a, false)
	if !slice.IsSorted(a, false) {
		t.Fail()
	}
}
