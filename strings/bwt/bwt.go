package bwt

// BWT - burrows-wheeler transform
func BWT(input string) (string, int) {
	n := len(input)
	s := input + input
	suffixes := make([]int, n)
	for i := 0; i < n; i++ {
		suffixes[i] = i
	}
	// sort suffixes
	less := func(i, j int) bool {
		for k := 0; k < n; k++ {
			if s[suffixes[i]+k] != s[suffixes[j]+k] {
				return s[suffixes[i]+k] < s[suffixes[j]+k]
			}
		}
		return false
	}
	// simple insertion sort
	for i := 1; i < n; i++ {
		j := i
		for j > 0 && less(j, j-1) {
			suffixes[j], suffixes[j-1] = suffixes[j-1], suffixes[j]
			j--
		}
	}
	// build BWT result
	result := make([]byte, n)
	originalIndex := 0
	for i := 0; i < n; i++ {
		if suffixes[i] == 0 {
			originalIndex = i
		}
		result[i] = s[suffixes[i]+n-1]
	}
	return string(result), originalIndex
}

// InverseBWT - inverse burrows-wheeler transform
func InverseBWT(bwt string, originalIndex int) string {
	n := len(bwt)
	count := make(map[byte]int)
	for i := 0; i < n; i++ {
		count[bwt[i]]++
	}
	// build first column
	firstCol := make([]byte, n)
	sortedKeys := make([]byte, 0, len(count))
	for k := range count {
		sortedKeys = append(sortedKeys, k)
	}
	// simple insertion sort
	for i := 1; i < len(sortedKeys); i++ {
		j := i
		for j > 0 && sortedKeys[j] < sortedKeys[j-1] {
			sortedKeys[j], sortedKeys[j-1] = sortedKeys[j-1], sortedKeys[j]
			j--
		}
	}
	idx := 0
	for _, k := range sortedKeys {
		c := count[k]
		for i := 0; i < c; i++ {
			firstCol[idx] = k
			idx++
		}
	}
	// build next array
	next := make([]int, n)
	tally := make(map[byte]int)
	for i := 0; i < n; i++ {
		c := bwt[i]
		next[tally[c]] = i
		tally[c]++
	}
	// reconstruct original string
	result := make([]byte, n)
	idx = originalIndex
	for i := 0; i < n; i++ {
		result[n-1-i] = bwt[idx]
		idx = next[idx]
	}
	return string(result)
}
