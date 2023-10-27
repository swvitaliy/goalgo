package goalgo

import (
	"fmt"
	gglCmp "github.com/google/go-cmp/cmp"
	"testing"
)

func TestLcs1(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	b := []int{0, 1, 3, 4, 5, 6, 7}
	c := Lcs(a, b)
	if !gglCmp.Equal(c, []int{1, 3, 7}) {
		fmt.Printf("expected: %v, got: %v\n", []int{1, 3, 7}, c)
		t.Fail()
	}
}
