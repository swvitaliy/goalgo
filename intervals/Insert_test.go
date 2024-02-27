package intervals

import (
	"reflect"
	"testing"
)

// Given a non-empty list of intervals, when inserting a new interval that does not overlap with any existing interval, then the resulting list should contain all the original intervals plus the new interval.
func TestInsertNoOverlap(t *testing.T) {
	a := [][]int{{1, 3}, {6, 7}, {9, 11}}
	x := []int{4, 5}
	expected := [][]int{{1, 3}, {4, 5}, {6, 7}, {9, 11}}
	result := Insert(a, x)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given a non-empty list of intervals, when inserting a new interval that overlaps with one existing interval, then the resulting list should contain all the original intervals merged with the new interval.
func TestInsertOverlapOne(t *testing.T) {
	a := [][]int{{1, 3}, {5, 7}, {9, 11}}
	x := []int{4, 8}
	expected := [][]int{{1, 3}, {4, 8}, {9, 11}}
	result := Insert(a, x)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given a non-empty list of intervals, when inserting a new interval that overlaps with multiple existing intervals, then the resulting list should contain all the original intervals merged with the new interval and any overlapping intervals.
func TestInsertOverlapMultiple(t *testing.T) {
	a := [][]int{{1, 3}, {5, 7}, {9, 11}}
	x := []int{4, 10}
	expected := [][]int{{1, 3}, {4, 11}}
	result := Insert(a, x)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given a non-empty list of intervals, when inserting a new interval that is completely contained within an existing interval, then the resulting list should contain all the original intervals and not the new interval.
func TestInsertContained(t *testing.T) {
	a := [][]int{{1, 3}, {5, 7}, {9, 11}}
	x := []int{6, 7}
	expected := [][]int{{1, 3}, {5, 7}, {9, 11}}
	result := Insert(a, x)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given a non-empty list of intervals, when inserting a new interval that completely contains an existing interval, then the resulting list should contain all the original intervals merged with the new interval.
func TestInsertContains(t *testing.T) {
	a := [][]int{{1, 3}, {5, 7}, {9, 11}}
	x := []int{4, 12}
	expected := [][]int{{1, 3}, {4, 12}}
	result := Insert(a, x)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Given a non-empty list of intervals, when inserting a new interval that is adjacent to an existing interval, then the resulting list should contain all the original intervals and the new interval.
func TestInsertAdjacent(t *testing.T) {
	a := [][]int{{1, 3}, {5, 7}, {9, 11}}
	x := []int{4, 5}
	expected := [][]int{{1, 3}, {4, 7}, {9, 11}}
	result := Insert(a, x)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
