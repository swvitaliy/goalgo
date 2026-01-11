package skiplist

import "testing"

func FuzzSkipList_InsertSearch(f *testing.F) {
	f.Add([]byte{
		0x01, 0x02, 0x03, 0x01, // key
		0x02, 0x03, 0x01, 0x02, // value
	})

	f.Fuzz(func(t *testing.T, data []byte) {
		sl := NewSkipList[int, int](RandomGenerator)
		m := make(map[int]int)

		// каждые 8 байт используем как пару int32 key/value
		for i := 0; i+7 < len(data); i += 8 {
			key := int(int32(data[i]) | int32(data[i+1])<<8 | int32(data[i+2])<<16 | int32(data[i+3])<<24)
			value := int(int32(data[i+4]) | int32(data[i+5])<<8 | int32(data[i+6])<<16 | int32(data[i+7])<<24)

			m[key] = value
			sl.Insert(key, value)
		}

		// проверка всех вставленных элементов
		for k, v := range m {
			if got, ok := sl.Search(k); !ok || got != v {
				t.Fatalf("mismatch key=%d: got=%d, want=%d", k, got, v)
			}
		}
	})
}
