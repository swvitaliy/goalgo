package gtyp

import (
	"github.com/google/go-cmp/cmp"
	"sort"
	"testing"
)

func TestTSliceC1(t *testing.T) {
	var s = TSliceC[int]{
		[]int{2, 5, 8, 1, 0, 3},
		func(a, b int) bool {
			return a < b
		},
	}
	sort.Sort(s)

	b := []int{0, 1, 2, 3, 5, 8}
	if cmp.Equal(b, s.a) {
		// OK
	} else {
		t.Fail()
	}
}
