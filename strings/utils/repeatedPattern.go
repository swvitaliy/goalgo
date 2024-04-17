package utils

func RepeatedPattern(s string) (string, bool) {
	l := len(s)
	for pl := 1; pl <= l/2; pl++ {
		if l%pl != 0 {
			continue
		}
		i := pl
		for ; i < len(s); i += pl {
			if s[0:pl] != s[i:i+pl] {
				i = -1
				break
			}
		}
		if i != -1 {
			return s[0:pl], true
		}
	}
	return "", false
}
