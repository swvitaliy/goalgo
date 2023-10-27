package list

func FromSlice[T any](a []T) *Node[T] {
	if len(a) == 0 {
		return nil
	}
	return &Node[T]{
		Val:  a[0],
		Next: FromSlice(a[1:]),
	}
}

func ToSlice[T any](node *Node[T]) []T {
	ans := make([]T, 0)
	for node != nil {
		ans = append(ans, node.Val)
		node = node.Next
	}
	return ans
}
