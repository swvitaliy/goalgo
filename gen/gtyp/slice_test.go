package gtyp

import (
	"github.com/google/go-cmp/cmp"
	"sort"
	"testing"
)

func TestTSlice1(t *testing.T) {
	var s = TSlice[int, LessOrdered[int]]{2, 5, 8, 1, 0, 3}
	sort.Sort(s)

	b := TSlice[int, LessOrdered[int]]{0, 1, 2, 3, 5, 8}
	if cmp.Equal(b, s) {
		// OK
	} else {
		t.Fail()
	}
}
