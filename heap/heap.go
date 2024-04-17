package heap

type Heap[T any] interface {
	Len() int
	IsEmpty() bool
	Peek() (T, bool)
	Enqueue(v T)
	Dequeue() (T, bool)
}
