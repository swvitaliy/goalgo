package algo

import (
	"github.com/google/go-cmp/cmp"
	"sort"
	"testing"
)

func TestTSliceC1(t *testing.T) {
	var s = TSlice[int]{
		[]int{2, 5, 8, 1, 0, 3},
		Less[int],
	}
	sort.Sort(s)

	b := []int{0, 1, 2, 3, 5, 8}
	if cmp.Equal(b, s.Data()) {
		// OK
	} else {
		t.Fail()
	}
}
