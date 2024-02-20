package aho_corasick

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindAll1(t *testing.T) {
	patterns := []string{"a", "ab", "abc", "abcd"}
	matches := SearchAll(patterns, "abcd")
	assert.Equal(t, []Match{{0, 0}, {0, 1}, {0, 2}, {0, 3}}, matches)
}

func TestFindAll2(t *testing.T) {
	patterns := []string{"a", "ab", "abc", "abcd"}
	matches := SearchAll(patterns, "abc")
	assert.Equal(t, []Match{{0, 0}, {0, 1}, {0, 2}}, matches)
}

func TestFindAll3(t *testing.T) {
	patterns := []string{"a", "ab", "abc", "abcd"}
	matches := SearchAll(patterns, "")
	assert.Equal(t, []Match{}, matches)
}

func TestFindAll4(t *testing.T) {
	patterns := []string{"ello", "worl", "world", "hekko"}
	matches := SearchAll(patterns, "hello world")
	assert.Equal(t, []Match{{1, 0}, {6, 1}, {6, 2}}, matches)
}

func TestFindAll5(t *testing.T) {
	patterns := []string{"hello", "world"}
	matches := SearchAll(patterns, "hello world")
	assert.Equal(t, []Match{{0, 0}, {6, 1}}, matches)
}
