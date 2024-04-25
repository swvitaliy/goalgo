package subseq

import (
	"fmt"
	gglCmp "github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLis1(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	b := Lis(a)
	assert.Equal(t, 3, len(b), "Lis failed")
	fmt.Printf("%v\n", b)
	if !gglCmp.Equal(b, []int{1, 3, 7}) && !gglCmp.Equal(b, []int{1, 2, 7}) {
		t.Fail()
	}
}

func TestLisLen1(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	assert.Equal(t, 3, LisLen(a), "LisLen failed")

}
