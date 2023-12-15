package kmp

// Knuth-Morris-Pratt algorithm

func prefixKmp(s string) []int {
	n := len(s)
	p := make([]int, n)
	p[0] = 0
	for i := 1; i < n; i++ {
		j := p[i-1]
		for j > 0 && s[i] != s[j] {
			j = p[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		p[i] = j
	}

	return p
}

func SearchAll(p, s string) []int {
	pi := prefixKmp(p + string(rune(0)) + s)
	ans := make([]int, 0)
	for i := range pi {
		if pi[i] == len(p) {
			ans = append(ans, i-len(p)-1)
		}
	}
	return ans
}

func SearchCount(p, s string) int {
	pi := prefixKmp(p + string(rune(0)) + s)
	ans := 0
	for i := range pi {
		if pi[i] == len(p) {
			ans++
		}
	}
	return ans
}

func SearchFirst(p, s string) int {
	pi := prefixKmp(p + string(rune(0)) + s)
	for i := range pi {
		if pi[i] == len(p) {
			return i - len(p) - 1
		}
	}
	return -1
}
