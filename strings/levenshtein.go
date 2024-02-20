package strings

type LevenshteinCostsConfig struct {
	Insert, Delete, Replace int
}

var defaultCosts = LevenshteinCostsConfig{
	Insert:  1,
	Replace: 1,
	Delete:  1,
}

func Levenshtein(s, t string) int {
	return LevenshteinWithCosts(s, t, defaultCosts)
}

func LevenshteinWithCosts(s, t string, cost LevenshteinCostsConfig) int {
	if len(s) == 0 {
		return len(t)
	}
	if len(t) == 0 {
		return len(s)
	}
	d := make([][]int, 2)
	d[0] = make([]int, len(t)+1)
	d[1] = make([]int, len(t)+1)
	d[0][0] = 0
	for j := 1; j < len(t)+1; j++ {
		d[0][j] = d[0][j-1] + cost.Insert
	}
	for i := 1; i < len(s)+1; i++ {
		cur := i & 1
		prev := 1 - cur
		d[cur][0] = d[prev][0] + cost.Delete
		for j := 1; j < len(t)+1; j++ {
			if s[i-1] == t[j-1] {
				d[cur][j] = d[prev][j-1]
			} else {
				d[cur][j] = min(
					d[prev][j]+cost.Delete,
					d[cur][j-1]+cost.Insert,
					d[prev][j-1]+cost.Replace,
				)
			}
		}
	}
	return d[len(s)&1][len(t)]
}
