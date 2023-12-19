package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevenshtein(t *testing.T) {
	assert.Equal(t, 0, Levenshtein("", ""))
	assert.Equal(t, 1, Levenshtein("a", ""))
	assert.Equal(t, 1, Levenshtein("", "a"))
	assert.Equal(t, 1, Levenshtein("a", "b"))
	assert.Equal(t, 1, Levenshtein("b", "a"))
	assert.Equal(t, 1, Levenshtein("abc", "abd"))
	assert.Equal(t, 1, Levenshtein("abc", "abe"))
	assert.Equal(t, 0, Levenshtein("abc", "abc"))
	assert.Equal(t, 3, Levenshtein("kitten", "sitting"))
}
