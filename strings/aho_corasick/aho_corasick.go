package aho_corasick

// Vertex is a node in the trie
type Vertex struct {
	children map[rune]*Vertex
	suf      *Vertex
	output   []int
}

type Searcher struct {
	root *Vertex
}

func (a *Searcher) AddString(s string, k int) *Vertex {
	n := a.root
	for _, c := range s {
		t, ok := n.children[c]
		if !ok {
			t = &Vertex{
				children: make(map[rune]*Vertex),
				suf:      nil,
				output:   nil,
			}
			n.children[c] = t
		}
		n = t
	}
	n.output = []int{k}
	return n
}

func NewSearcher(patterns []string) *Searcher {
	root := &Vertex{
		children: make(map[rune]*Vertex),
		suf:      nil,
		output:   nil,
	}

	a := Searcher{root}
	for i, p := range patterns {
		a.AddString(p, i)
	}

	a.PrepareSuffixes()
	return &a
}

func SearchFirst(patterns []string, text string) (int, []int) {
	s := NewSearcher(patterns)
	return s.searchFirst(patterns, text)
}

func (a *Searcher) searchFirst(patterns []string, text string) (int, []int) {
	var v *Vertex = a.root
	for i, c := range text {
		for v.children[c] == nil && v != a.root {
			v = v.suf
		}

		if v.children[c] != nil {
			v = v.children[c]
		}

		if v.output == nil {
			continue
		}

		if len(v.output) > 0 {
			return i - len(patterns[v.output[0]]) + 1, v.output
		}
	}

	return -1, nil
}

type Match struct {
	Pos    int
	Output int
}

func SearchAll(patterns []string, text string) []Match {
	s := NewSearcher(patterns)
	return s.searchAll(patterns, text)
}

func (a *Searcher) searchAll(patterns []string, text string) []Match {
	ans := make([]Match, 0)
	var v *Vertex = a.root
	for i, c := range text {
		for v.children[c] == nil && v != a.root {
			v = v.suf
		}

		if v.children[c] != nil {
			v = v.children[c]
		}

		if v.output == nil {
			continue
		}

		for _, p := range v.output {
			ans = append(ans, Match{i - len(patterns[p]) + 1, p})
		}
	}

	return ans
}

func (a *Searcher) PrepareSuffixes() {
	queue := make([]*Vertex, 0)
	for _, u := range a.root.children {
		u.suf = a.root
		queue = append(queue, u)
	}

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		for c, u := range v.children {
			nextSuf := v.suf
			for nextSuf != nil && nextSuf.children[c] == nil {
				nextSuf = nextSuf.suf
			}

			if nextSuf == nil {
				u.suf = a.root
			} else {
				u.suf = nextSuf.children[c]
			}

			u.output = append(u.output, u.suf.output...)
			queue = append(queue, u)
		}
	}
}
