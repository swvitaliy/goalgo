package rabinKarp

const p = 1000000007

var pp []int

func initPP(n int) {
	var k int
	if pp == nil {
		pp = make([]int, n)
		pp[0] = 1
		k = 1
	} else {
		k = len(pp)
	}

	if n <= k {
		return
	}

	for i := k; i < n; i++ {
		pp[i] = pp[i-1] * p
	}
}

// Search searches for s in t
func Search(s, t string) int {
	n := len(s)
	m := len(t)

	if n == 0 || m == 0 {
		return -1
	}

	if n > m {
		return -1
	}

	initPP(m)

	hs := 0
	for i, c := range s {
		hs += int(c) * pp[i]
	}

	ht := make([]int, m)
	for i, c := range t {
		ht[i] += int(c) * pp[n-1]
		if i > 0 {
			ht[i] += ht[i-1]
		}
	}

	for i := 0; i < m-n; i++ {
		curH := ht[i+n-1]
		if i > 0 {
			curH -= ht[i-1]
		}
		if hs*pp[i] == curH {
			return i
		}
	}

	return -1
}
