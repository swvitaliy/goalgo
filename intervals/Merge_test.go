package intervals

import (
	"reflect"
	"testing"
)

// Given a 2D slice of intervals with no overlapping intervals, the function should return the same slice.
func TestNoOverlappingIntervals(t *testing.T) {
	input := [][]int{{1, 3}, {4, 6}, {7, 9}}
	expected := [][]int{{1, 3}, {4, 6}, {7, 9}}
	result := Merge(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given a 2D slice of intervals with overlapping intervals, the function should merge the intervals and return a new slice with the merged intervals.
func TestOverlappingIntervals(t *testing.T) {
	input := [][]int{{1, 3}, {2, 6}, {8, 10}}
	expected := [][]int{{1, 6}, {8, 10}}
	result := Merge(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given a 2D slice of intervals with some intervals completely contained within others, the function should merge the intervals and return a new slice with the merged intervals.
func TestContainedIntervals(t *testing.T) {
	input := [][]int{{1, 10}, {2, 6}, {4, 8}}
	expected := [][]int{{1, 10}}
	result := Merge(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given an empty 2D slice of intervals, the function should return an empty slice.
func TestEmptyIntervals(t *testing.T) {
	input := [][]int{}
	expected := [][]int{}
	result := Merge(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given a 2D slice of intervals with only one interval, the function should return the same slice.
func TestSingleInterval(t *testing.T) {
	input := [][]int{{1, 5}}
	expected := [][]int{{1, 5}}
	result := Merge(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given a 2D slice of intervals with intervals that have negative values, the function should merge the intervals and return a new slice with the merged intervals.
func TestNegativeIntervals(t *testing.T) {
	input := [][]int{{-5, -2}, {-4, 0}, {-3, -1}}
	expected := [][]int{{-5, 0}}
	result := Merge(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
