package radix

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRadixTrie(t *testing.T) {
	a := NewRadixTrie()
	assert.True(t, a != nil)
}

func TestRadixFewStrings2(t *testing.T) {
	a := NewRadixTrie()
	a.AddString("apple")
	assert.True(t, a.SearchString("apple"))
	assert.False(t, a.SearchString("app"))
	assert.True(t, a.SearchPrefix("app"))
	a.AddString("app")
	assert.True(t, a.SearchString("app"))
}
func TestRadixFewStrings(t *testing.T) {
	a := NewRadixTrie()
	a.AddString("helo")
	a.AddString("hel")
	a.AddString("hello")
	a.AddString("hello123")
	a.AddString("hllo")
	assert.True(t, a.SearchPrefix("hello"))
	assert.True(t, a.SearchString("hello"))
	assert.True(t, a.SearchString("hello123"))
	assert.False(t, a.SearchString("hello456"))
}
func TestRadixOneString(t *testing.T) {
	a := NewRadixTrie()
	a.AddString("hello")
	assert.True(t, a.SearchPrefix("hello"))
	assert.True(t, a.SearchString("hello"))
	assert.False(t, a.SearchString("hello123"))
	assert.False(t, a.SearchString("hel"))
}
