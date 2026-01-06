package skip_list

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkipList_RangeInts(t *testing.T) {
	tests := []struct {
		name     string
		skipList *SkipList[int, struct{}]
		begin    int
		end      int
		result   []int
	}{
		{
			name:     "empty skip list",
			skipList: FromSliceOfKeys([]int{}),
			begin:    0,
			end:      10,
			result:   []int{},
		},
		{
			name:     "simple range - range more left",
			skipList: FromSliceOfKeys([]int{1, 2, 3}),
			begin:    -1,
			end:      0,
			result:   []int{},
		},
		{
			name:     "simple range - range more right",
			skipList: FromSliceOfKeys([]int{1, 2, 3}),
			begin:    10,
			end:      12,
			result:   []int{},
		},
		{
			name:     "simple range - range exactly coinsident with skip list first/last keys",
			skipList: FromSliceOfKeys([]int{1, 2, 3}),
			begin:    1,
			end:      3,
			result:   []int{1, 2, 3},
		},
		{
			name:     "unordered skip list",
			skipList: FromSliceOfKeys([]int{7, 0, 12, 5, 63, 24, 76}),
			begin:    5,
			end:      65,
			result:   []int{5, 7, 12, 24, 63},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []int
			for x := range tt.skipList.Range(tt.begin, tt.end) {
				result = append(result, x)
			}
			require.True(t, slices.Equal(result, tt.result))

		})
	}
}
