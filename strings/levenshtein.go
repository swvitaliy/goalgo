package strings

const (
	levenshteinDefaultCost = 1
	levenshteinInsertCost  = 1
	levenshteinDeleteCost  = 1
	levenshteinReplaceCost = 1
)

func Levenshtein(s, t string) int {
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
		d[0][j] = d[0][j-1] + levenshteinInsertCost
	}
	for i := 1; i < len(s)+1; i++ {
		cur := i & 1
		prev := 1 - cur
		d[cur][0] = d[prev][0] + levenshteinDeleteCost
		for j := 1; j < len(t)+1; j++ {
			if s[i-1] == t[j-1] {
				d[cur][j] = d[prev][j-1]
			} else {
				d[cur][j] = min(
					d[prev][j]+levenshteinDeleteCost,
					d[cur][j-1]+levenshteinInsertCost,
					d[prev][j-1]+levenshteinReplaceCost,
				)
			}
		}
	}
	return d[len(s)&1][len(t)]
}
