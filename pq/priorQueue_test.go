package pq

import "testing"

func TestPQueue_Enqueue_Dequeue(t *testing.T) {
	q := NewPQ[uint]()
	q.Enqueue(0, 0)
	q.Enqueue(7, 7)
	q.Enqueue(2, 22)
	if v, _ := q.Dequeue(); v != 0 {
		t.Error()
	}
	if v, _ := q.Dequeue(); v != 22 {
		t.Error()
	}
	if v, _ := q.Dequeue(); v != 7 {
		t.Error()
	}
}
