package intervals

/*
func Intersect2(a [][]int, b [][]int) [][]int {
  c := make([][]int, len(a)+len(b))
  copy(c, a)
  copy(c[len(a):],b)
  sort.Slice(c, func(i, j int) bool {
    return c[i][0] < c[j][0]
  })
  return Intersect(c)
}

*/

func Intersect(a [][]int) [][]int {
	c := make([][]int, 0, len(a))
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a) && a[j][0] <= a[i][1]; j++ {
			c = append(c, []int{a[j][0], min(a[i][1], a[j][1])})
		}
	}
	return c
}

func Intersect2(a [][]int, b [][]int) [][]int {
	i := 0
	j := 0
	c := make([][]int, 0, len(a)+len(b))

	for i < len(a) && j < len(b) {
		ai := a[i]
		bj := b[j]

		s := max(ai[0], bj[0])
		f := min(ai[1], bj[1])
		if s <= f {
			c = append(c, []int{s, f})
		}

		if ai[1] < bj[1] {
			i++
		} else {
			j++
		}
	}

	return c
}
