package trie

type RadixNode struct {
	eow   bool
	edges map[byte]*RadixEdge
}

type RadixEdge struct {
	chars []byte
	next  *RadixNode
}

type RadixTrie struct {
	root *RadixNode
}

func NewRadixTrie() *RadixTrie {
	return &RadixTrie{
		root: &RadixNode{false, make(map[byte]*RadixEdge)},
	}
}

func (t *RadixTrie) AddString(s string) {
	n := t.root
	w := []byte(s)
	for len(w) > 0 {
		c := w[0]
		if n.edges == nil {
			n.edges = make(map[byte]*RadixEdge)
		}
		e, ok := n.edges[c]
		if !ok {
			n.edges[c] = &RadixEdge{
				chars: w,
				next:  &RadixNode{eow: true},
			}
			break
		}

		pw := w
		m := len(e.chars)
		if m < len(w) {
			pw = w[:m]
		}
		j := findFirstDiff(pw, e.chars)
		if j != m {
			tail := e.chars[j:]
			e.chars = e.chars[:j]
			tailNode := &RadixNode{false, make(map[byte]*RadixEdge)}
			tailNode.edges[tail[0]] = &RadixEdge{
				chars: tail,
				next:  e.next,
			}
			e.next = tailNode
		}

		if len(w) == len(e.chars) {
			e.next.eow = true
		}

		n = e.next
		w = w[j:]
	}
}

func findFirstDiff(a, b []byte) int {
	if len(a) > len(b) {
		a, b = b, a
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return len(a)
}

func (t *RadixTrie) searchRadixNode(s string) *RadixNode {
	n := t.root
	w := []byte(s)
	for len(w) > 0 {
		c := w[0]
		e, ok := n.edges[c]
		if !ok {
			return nil
		}
		j := findFirstDiff(w, e.chars)
		if j != len(e.chars) {
			return nil
		}
		n = e.next
		w = w[j:]
	}
	return n
}

func (t *RadixTrie) DeleteString(s string) bool {
	n := t.searchRadixNode(s)
	if n == nil {
		return false
	}
	if n.eow {
		n.eow = false
		return true
	}
	return false
}

func (t *RadixTrie) SearchString(s string) bool {
	n := t.searchRadixNode(s)
	if n == nil {
		return false
	}
	return n.eow
}

func (t *RadixTrie) SearchPrefix(s string) bool {
	n := t.root
	w := []byte(s)
	i := 0
	for len(w) > 0 {
		c := w[0]
		e, ok := n.edges[c]
		if !ok {
			break
		}
		j := findFirstDiff(w, e.chars)
		if j != len(e.chars) {
			w = w[j:]
			i += j
			break
		}
		n = e.next
		w = w[j:]
		i += j
	}
	return len(w) == 0
}
