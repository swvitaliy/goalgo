package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatedPattern1(t *testing.T) {
	s := "abababababababababab"
	p, ok := RepeatedPattern(s)
	expected := "ab"
	assert.True(t, ok)
	assert.Equal(t, expected, p)
}

func TestRepeatedPattern2(t *testing.T) {
	s := "ababababababababababa"
	_, ok := RepeatedPattern(s)
	assert.False(t, ok)
}
