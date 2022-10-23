package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"lis/typ"
	"testing"
)

func TestLis1(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	b := LisNumbers(a, typ.IntBounds)

	for _, v := range b {
		fmt.Print(v, " ")
	}
	fmt.Println()

	if cmp.Equal(b, []int{1, 3, 7}) || cmp.Equal(b, []int{1, 2, 7}) {
		// OK
	} else {
		t.Fail()
	}
}
