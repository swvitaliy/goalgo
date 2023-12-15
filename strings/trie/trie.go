package trie

type TrieNode struct {
	children map[byte]*TrieNode
}

func NewTrie() *TrieNode {
	return &TrieNode{make(map[byte]*TrieNode)}
}

func (n *TrieNode) AddString(s string) *TrieNode {
	for _, c := range []byte(s) {
		t, ok := n.children[c]
		if !ok {
			t = &TrieNode{make(map[byte]*TrieNode)}
			n.children[c] = t
		}
		n = t
	}
	n.children['$'] = &TrieNode{}
	return n
}

func (n *TrieNode) SearchString(s string) bool {
	_, ok := n.SearchStringNode(s)
	return ok
}

func (n *TrieNode) SearchStringNode(s string) (*TrieNode, bool) {
	n, hasPrefix := n.SearchPrefixNode(s)
	if !hasPrefix {
		return nil, false
	}
	_, ok := n.children['$']
	return n, ok
}

func (n *TrieNode) SearchPrefix(s string) bool {
	_, ok := n.SearchPrefixNode(s)
	return ok
}

func (n *TrieNode) SearchPrefixNode(s string) (*TrieNode, bool) {
	for _, c := range []byte(s) {
		t, ok := n.children[c]
		if !ok {
			return n, false
		}
		n = t
	}
	return n, true

}

func (n *TrieNode) DeleteString(s string) {
	if n, ok := n.SearchStringNode(s); ok {
		delete(n.children, '$')
	}
}
