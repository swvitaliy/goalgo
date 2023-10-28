package goalgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLcs1(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	b := []int{0, 1, 3, 4, 5, 6, 7}
	assert.Equal(t, 3, LcsLen(a, b), "LcsLen failed")
	c := Lcs(a, b)
	assert.Equal(t, []int{1, 3, 7}, c, "Lcs failed")
}
