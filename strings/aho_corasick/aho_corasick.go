package aho_corasick

// Vertex is a node in the trie
type Vertex struct {
	children map[byte]*Vertex
	suf      *Vertex
}

type Searcher struct {
	root  *Vertex
	tnMap map[*Vertex]int // terminal nodes map to pattern indexes
}

type Match struct {
	Pos          int
	PatternIndex int
}

func (n *Vertex) AddString(s string) *Vertex {
	for _, c := range []byte(s) {
		t, ok := n.children[c]
		if !ok {
			t = &Vertex{make(map[byte]*Vertex), nil}
			n.children[c] = t
		}
		n = t
	}
	n.children['$'] = &Vertex{}
	return n
}

func NewSearcher(patterns []string) *Searcher {
	root := &Vertex{make(map[byte]*Vertex), nil}
	tnMap := make(map[*Vertex]int)
	for i, p := range patterns {
		termNode := root.AddString(p)
		tnMap[termNode] = i
	}
	PrepareSuffixes(root)

	return &Searcher{root, tnMap}
}

func SearchFirst(patterns []string, text string) Match {
	s := NewSearcher(patterns)
	return s.SearchFirst(text)
}

func (s *Searcher) SearchFirst(text string) Match {
	var v *Vertex = s.root
	for i, c := range []byte(text) {
		for ; v.children[c] == nil; v = v.suf {
		}
		v = v.children[c]

		if v.children['$'] != nil {
			continue
		}

		if v == s.root {
			continue
		}

		return Match{i, s.tnMap[v]}
	}

	return Match{-1, -1}
}

func SearchCount(patterns []string, text string) map[int]int {
	s := NewSearcher(patterns)
	return s.SearchCount(text)
}

func (s *Searcher) SearchCount(text string) map[int]int {
	root := s.root
	tnMap := s.tnMap

	ans := make(map[int]int)
	var v *Vertex = root
	for _, c := range []byte(text) {
		for ; v.children[c] == nil; v = v.suf {
		}
		v = v.children[c]

		if v.children['$'] != nil {
			continue
		}

		if v == root {
			continue
		}

		ans[tnMap[v]]++
	}

	return ans
}

func SearchAll(patterns []string, text string) []Match {
	s := NewSearcher(patterns)
	return s.SearchAll(text)
}

func (s *Searcher) SearchAll(text string) []Match {
	root := s.root
	tnMap := s.tnMap

	ans := make([]Match, 0)
	var v *Vertex = root
	for i, c := range []byte(text) {
		for ; v.children[c] == nil; v = v.suf {
		}
		v = v.children[c]

		if v.children['$'] != nil {
			continue
		}

		if v == root {
			continue
		}

		ans = append(ans, Match{i, tnMap[v]})
	}

	return ans
}

func PrepareSuffixes(root *Vertex) {
	bfs(root)
}

func bfs(root *Vertex) []*Vertex {
	queue := []*Vertex{root}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for parentChar, child := range v.children {
			parentSuf := v.suf
			nextSuf := parentSuf.children[parentChar].suf
			child.suf = nextSuf
			queue = append(queue, child)
		}
	}
	return queue
}
