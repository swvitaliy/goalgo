package bloom

import (
	"github.com/spaolacci/murmur3"
	"hash"
)

type BF struct {
	m      []uint64
	size   uint64
	hashes []hash.Hash64
	k      uint32
}

func NewBloomFilter(n uint64, k uint32) *BF {
	hashes := make([]hash.Hash64, k)
	for i := uint32(0); i < k; i++ {
		hashes[i] = murmur3.New64WithSeed(i)
	}
	return &BF{
		m:      make([]uint64, n/64+1),
		size:   n,
		hashes: hashes,
		k:      k,
	}
}

func (bf *BF) Add(data []byte) {
	for _, h := range bf.hashes {
		h.Reset()
		writes, err := h.Write(data)
		if err != nil {
			panic("error writing to hash")
		}
		if writes != len(data) {
			panic("hash wrote less bytes than expected")
		}
		i := h.Sum64() % bf.size
		bf.m[i/64] |= 1 << uint(i%64)
	}
}

func (bf *BF) AddString(x string) {
	bf.Add([]byte(x))
}

func (bf *BF) Test(data []byte) bool {
	for _, h := range bf.hashes {
		h.Reset()
		writes, err := h.Write(data)
		if err != nil {
			panic("error writing to hash")
		}
		if writes != len(data) {
			panic("hash wrote less bytes than expected")
		}
		i := h.Sum64() % bf.size
		if bf.m[i/64]&(1<<uint(i%64)) == 0 {
			return false
		}
	}
	return true
}

func (bf *BF) TestString(x string) bool {
	return bf.Test([]byte(x))
}
