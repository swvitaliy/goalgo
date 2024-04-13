package trie

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTrie(t *testing.T) {
	a := NewTrieNode()
	assert.True(t, a != nil)
}

func TestAddString(t *testing.T) {
	a := NewTrieNode()
	a.AddString("hello")
	assert.True(t, a.SearchPrefix("hello"))
	assert.True(t, a.SearchString("hello"))
	assert.False(t, a.SearchString("hello123"))
	assert.False(t, a.SearchString("hel"))
	assert.False(t, a.SearchString("helo"))

	a.DeleteString("hello")
	assert.False(t, a.SearchString("hello"))

}
