package misc

// MajorityElement returns the majority element in arr.
// Majority element needs to be appearing more than ⌊n / 2⌋ times.
// https://en.wikipedia.org/wiki/Boyer%E2%80%93Moore_majority_vote_algorithm
func MajorityElement[T comparable](s []T) T {
	ans, f := s[0], 1
	for i := 1; i < len(s); i++ {
		if f == 0 {
			ans, f = s[i], 1
		} else {
			if s[i] == ans {
				f++
			} else {
				f--
			}
		}
	}
	return ans
}

func IsMajorityElement[T comparable](s []T, v T) bool {
	f := 0
	for _, x := range s {
		if x == v {
			f++
		}
	}
	return f > len(s)/2
}
