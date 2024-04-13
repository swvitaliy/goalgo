package subseq

import (
	"goalgo/slices"
)

func Lcs[S ~[]T, T comparable](a, b S) S {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := 0; i <= n; i++ {
		dp[i][0] = 0
	}
	for j := 0; j <= m; j++ {
		dp[0][j] = 0
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	ans := make([]T, 0)
	i, j := n, m
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			ans = append(ans, a[i-1])
			i--
			j--
		} else if dp[i][j] == dp[i-1][j] {
			i--
		} else {
			j--
		}
	}

	slices.Reverse(ans)
	return ans
}

// LcsLen returns the length of the longest common subsequence of a and b
func LcsLen[S ~[]T, T comparable](a, b S) int {
	if len(a) < len(b) {
		a, b = b, a
	}
	m := len(a)
	n := len(b)
	// In i-th position of the c array, the i-th element is the length of the longest common subsequence of prefixes a[0:i] and b[0:i]
	c := make([]int, n+1)
	d := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			t := c[j+1]
			if a[i] == b[j] {
				c[j+1] = d + 1
			} else {
				if c[j+1] < c[j] {
					c[j+1] = c[j]
				}
			}
			d = t
		}
	}

	return c[n]
}
