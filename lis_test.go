package goalgo

import (
	"fmt"
	gglCmp "github.com/google/go-cmp/cmp"
	"testing"
)

func TestLis1(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	b := Lis(a)
	fmt.Printf("%v\n", b)
	if !gglCmp.Equal(b, []int{1, 3, 7}) && !gglCmp.Equal(b, []int{1, 2, 7}) {
		t.Fail()
	}
}
