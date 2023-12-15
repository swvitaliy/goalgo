package min_path

import (
	gglCmp "github.com/google/go-cmp/cmp"
	"testing"
)

func TestDijkstra1(t *testing.T) {
	e := make([][]edge, 5)
	e[0] = []edge{
		{1, 7},
		{2, 1},
	}
	e[1] = []edge{
		{2, 5},
		{3, 2},
	}
	e[2] = []edge{
		{3, 7},
	}
	w, path := Dijkstra(e, 0, 3)
	if w != 8 {
		t.Error()
	}
	if !gglCmp.Equal(path, []int{2, 3}) {
		t.Error()
	}
}
