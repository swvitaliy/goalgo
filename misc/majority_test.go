package misc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMajority1(t *testing.T) {
	a := []int{1, 1, 1, 2, 2, 3, 3, 2, 2, 2, 3, 2, 2}
	v := MajorityElement(a)
	assert.Equal(t, 2, v)
	assert.True(t, IsMajorityElement(a, 2))
	assert.False(t, IsMajorityElement(a, 1))
	assert.False(t, IsMajorityElement(a, 3))
}
