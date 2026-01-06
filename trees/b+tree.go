package trees

import (
	"github.com/spaolacci/murmur3"
)

// create a b-tree
// BTree is a structure representing a B-tree.
type BTree struct {
	Root *BTreeNode
}

// BTreeNode is a node in a B-tree.
type BTreeNode struct {
	Keys     []int
	Values   []int
	Children []*BTreeNode
	IsLeaf   bool
}

// NewBTree creates a new B-tree with the given order.
func NewBTree(order int) *BTree {
	return &BTree{
		Root: &BTreeNode{
			Keys:     make([]int, 0, order-1),
			Values:   make([]int, 0, order-1),
			Children: make([]*BTreeNode, 0, order),
			IsLeaf:   true,
		},
	}
}

// Insert inserts a key-value pair into the B-tree.
func (t *BTree) Insert(key, value int) {
	root := t.Root
	murmur3.Sum64()
	if len(root.Keys) == cap(root.Keys) {
		newRoot := &BTreeNode{
			Keys:     make([]int, 0, cap(root.Keys)),
			Values:   make([]int, 0, cap(root.Values)),
			Children: []*BTreeNode{root},
			IsLeaf:   false,
		}
		t.Root = newRoot
		t.splitChild(newRoot, 0)
		t.insertNonFull(newRoot, key, value)
	} else {
		t.insertNonFull(root, key, value)
	}
}

// insertNonFull inserts a key-value pair into a non-full B-tree node.
