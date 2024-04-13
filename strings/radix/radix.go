package radix

type Node struct {
	eow   bool
	edges map[rune]*Edge
}

type Edge struct {
	chars []rune
	next  *Node
}

type Trie struct {
	root *Node
}

func NewRadixTrie() *Trie {
	return &Trie{
		root: &Node{false, make(map[rune]*Edge)},
	}
}

func (t *Trie) AddString(s string) {
	n := t.root
	w := []rune(s)
	for len(w) > 0 {
		c := w[0]
		if n.edges == nil {
			n.edges = make(map[rune]*Edge)
		}
		e, ok := n.edges[c]
		if !ok {
			n.edges[c] = &Edge{
				chars: w,
				next:  &Node{eow: true},
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
			tailNode := &Node{false, make(map[rune]*Edge)}
			tailNode.edges[tail[0]] = &Edge{
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

func findFirstDiff(a, b []rune) int {
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

func (t *Trie) searchRadixNode(s string) *Node {
	n := t.root
	w := []rune(s)
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

func (t *Trie) DeleteString(s string) bool {
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

func (t *Trie) SearchString(s string) bool {
	n := t.searchRadixNode(s)
	if n == nil {
		return false
	}
	return n.eow
}

func (t *Trie) SearchPrefix(s string) bool {
	n := t.root
	w := []rune(s)
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
