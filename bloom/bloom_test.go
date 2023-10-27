package bloom

import "testing"

func TestBloomFilter_Add(t *testing.T) {
	a := NewBloomFilter(100, 10)
	a.AddString("ase fasdfahjs bfau shdfb")
	a.AddString("hello")
	if !a.TestString("hello") {
		t.Error("add_test failed")
	}
}
