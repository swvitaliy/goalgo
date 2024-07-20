package rabinKarp

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRabinKarpSearch0(t *testing.T) {
	var a int
	a = Search("hello world", "h3llo")
	assert.Equal(t, -1, a)
	a = Search("", "h3llo")
	assert.Equal(t, -1, a)
	a = Search("hello world", "")
	assert.Equal(t, -1, a)
	a = Search("hello world", "h3llo 11111111")
	assert.Equal(t, -1, a)
}

func TestRabinKarpSearch1(t *testing.T) {
	a := Search("hello world", "hello")
	assert.Equal(t, 0, a)
}

func TestRabinKarpSearch2(t *testing.T) {
	a := Search("hello world", "world")
	assert.Equal(t, 6, a)
}

func TestRabinKarpSearch3(t *testing.T) {
	a := Search("hello world", "h3llo")
	assert.Equal(t, -1, a)
}

func TestRabinKarpSearch4(t *testing.T) {
	a := Search("hello world", "h3llo 11111111")
	assert.Equal(t, -1, a)
}

func TestRabinKarpSearch5(t *testing.T) {
	a := Search("hello world", "hello world")
	assert.Equal(t, 0, a)
}

// It returns the correct index when the pattern is found in the text.
func TestPatternFound(t *testing.T) {
	s := "abc"
	text := "abcdefg"
	expected := 0

	result := Search(s, text)

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

// It returns -1 when the pattern is not found in the text.
func TestPatternNotFound(t *testing.T) {
	s := "xyz"
	text := "abcdefg"
	expected := -1

	result := Search(s, text)

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

// It returns -1 when the pattern is longer than the text.
func TestPatternLongerThanText(t *testing.T) {
	s := "abcdefg"
	text := "abc"
	expected := -1

	result := Search(s, text)

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

// It correctly handles patterns and texts with the maximum allowed length.
func TestMaxAllowedLength(t *testing.T) {
	s := "a" + strings.Repeat("b", 1000000)
	text := "c" + strings.Repeat("d", 1000000)
	expected := -1

	result := Search(s, text)

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

// It correctly handles patterns and texts with non-ASCII characters.
func TestNonASCIICharacters(t *testing.T) {
	s := "こんにちは"
	text := "こんにちは、世界"
	expected := 0

	result := Search(s, text)

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

// It correctly handles patterns and texts with special characters.
func TestSpecialCharacters(t *testing.T) {
	s := "!@#$%^&*()"
	text := "Hello!@#$%^&*()World"
	expected := 5

	result := Search(s, text)

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}
