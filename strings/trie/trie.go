package trie

type Node struct {
	children map[rune]*Node
	term     bool
}

func NewTrieNode() *Node {
	return &Node{
		children: make(map[rune]*Node),
		term:     false,
	}
}

func (n *Node) AddString(s string) *Node {
	for _, c := range []rune(s) {
		t, ok := n.children[c]
		if !ok {
			t = NewTrieNode()
			n.children[c] = t
		}
		n = t
	}
	n.term = true
	return n
}

func (n *Node) SearchString(s string) bool {
	_, ok := n.SearchStringNode(s)
	return ok
}

func (n *Node) SearchStringNode(s string) (*Node, bool) {
	n, hasPrefix := n.SearchPrefixNode(s)
	if !hasPrefix {
		return nil, false
	}
	return n, n.term
}

func (n *Node) SearchPrefix(s string) bool {
	_, ok := n.SearchPrefixNode(s)
	return ok
}

func (n *Node) SearchPrefixNode(s string) (*Node, bool) {
	for _, c := range []rune(s) {
		t, ok := n.children[c]
		if !ok {
			return n, false
		}
		n = t
	}
	return n, true

}

func (n *Node) DeleteString(s string) {
	if n, ok := n.SearchStringNode(s); ok {
		n.term = false
	}
}
