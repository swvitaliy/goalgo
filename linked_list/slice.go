package linked_list

func FromSlice[T any](a []T) *Node[T] {
	if len(a) == 0 {
		return nil
	}
	return &Node[T]{
		val:  a[0],
		next: FromSlice[T](a[1:]),
	}
}

func ToSlice[T any](node *Node[T]) []T {
	ans := make([]T, 0)
	for node != nil {
		ans = append(ans, node.val)
		node = node.next
	}
	return ans
}
