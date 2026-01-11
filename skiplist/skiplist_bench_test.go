package skiplist

import (
	"goalgo/slices/iter_utils"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

const (
	maxKey    = int64(10_000_000)
	listSize  = 5_000_000
	batchSize = 16
)

const concurrentParallelism = 4

const compactEvery = 10_000

const shuffleCoef = 0.3

// -------------------------------------------------------------
// Search benchmarks

func BenchmarkSkipList_Search(b *testing.B) {
	sl := NewFromKeysIter(
		RandomGenerator,
		iter_utils.RandomInt64Iter(listSize, 0, maxKey),
	)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
		for pb.Next() {
			sl.Search(r.Int63n(maxKey))
		}
	})
}

func BenchmarkSkipList_BatchSearchNodes(b *testing.B) {
	sl := NewFromKeysIter(
		RandomGenerator,
		iter_utils.RandomInt64Iter(listSize, 0, maxKey),
	)

	preparedKeys := make([][]int64, batchSize)

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	keys := make([]int64, batchSize)
	for j := 0; j < batchSize; j++ {
		for i := range keys {
			keys[i] = r.Int63n(maxKey)
		}
		preparedKeys[j] = keys
	}

	var batchCount atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = sl.BatchSearchNodes(preparedKeys[batchCount.Load()%batchSize])
			batchCount.Add(1)
		}
	})
}

func BenchmarkSkipList_BatchFriendlyLevelGenerator(b *testing.B) {
	sl := NewFromKeysIter(
		BatchFriendlyGenerator,
		iter_utils.RandomInt64Iter(listSize, 0, maxKey),
	)

	preparedKeys := make([][]int64, batchSize)

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	keys := make([]int64, batchSize)
	for j := 0; j < batchSize; j++ {
		for i := range keys {
			keys[i] = r.Int63n(maxKey)
		}
		preparedKeys[j] = iter_utils.PartiallyShuffledSlice(keys, shuffleCoef)
	}

	var batchCount atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = sl.BatchSearchNodes(preparedKeys[batchCount.Load()%batchSize])
			batchCount.Add(1)
		}
	})
}

func BenchmarkSkipList_ConcurrentSearch(b *testing.B) {
	sl := NewConcurrentFromKeysIter(
		RandomGenerator,
		iter_utils.RandomInt64Iter(listSize, 0, maxKey),
	)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(concurrentParallelism)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sl.Contains(rand.Int63n(maxKey))
		}
	})
}

// -------------------------------------------------------------
// Insert benchmarks

func BenchmarkSkipList_Insert(b *testing.B) {
	sl := NewSkipList[int64, struct{}](RandomGenerator)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
		for pb.Next() {
			sl.Insert(r.Int63n(maxKey), struct{}{})
		}
	})
}

func BenchmarkSkipList_ConcurrentInsert(b *testing.B) {
	sl := NewConcurrentSkipList[int64, struct{}](RandomGenerator)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(concurrentParallelism)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
		for pb.Next() {
			sl.Insert(r.Int63n(maxKey), struct{}{})
		}
	})
}

func BenchmarkSkipList_ConcurrentBulkInsert(b *testing.B) {
	sl := NewConcurrentFromKeysIter(
		RandomGenerator,
		iter_utils.RandomInt64Iter(listSize, 0, maxKey),
	)

	preparedKeys := make([][]int64, batchSize)
	preparedValues := make([][]struct{}, batchSize)

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	keys := make([]int64, batchSize)
	values := make([]struct{}, batchSize)
	for j := 0; j < batchSize; j++ {
		for i := range keys {
			keys[i] = r.Int63n(maxKey)
		}
		preparedKeys[j] = keys
		preparedValues[j] = values
	}

	var batchCount atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(concurrentParallelism)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sl.BulkInsert(
				preparedKeys[batchCount.Load()%batchSize],
				preparedValues[batchCount.Load()%batchSize],
			)
			batchCount.Add(1)
		}
	})
}

// -------------------------------------------------------------
// Delete benchmarks

func BenchmarkSkipList_Delete(b *testing.B) {
	sl := NewFromKeysIter(
		RandomGenerator,
		iter_utils.RandomInt64Iter(listSize, 0, maxKey),
	)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
		for pb.Next() {
			sl.Delete(r.Int63n(maxKey))
		}
	})
}

func BenchmarkSkipList_ConcurrentDelete(b *testing.B) {
	sl := NewConcurrentFromKeysIter(
		RandomGenerator,
		iter_utils.RandomInt64Iter(listSize, 0, maxKey),
	)

	var deleteCount atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(concurrentParallelism)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sl.Delete(rand.Int63n(maxKey))
			if deleteCount.Add(1)%compactEvery == 0 {
				sl.Compact()
			}
		}
	})
}

func BenchmarkSkipList_ConcurrentBulkDelete(b *testing.B) {
	sl := NewConcurrentFromKeysIter(
		BatchFriendlyGenerator,
		iter_utils.RandomInt64Iter(listSize, 0, maxKey),
	)

	preparedKeys := make([][]int64, batchSize)

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	keys := make([]int64, batchSize)
	for j := 0; j < batchSize; j++ {
		for i := range keys {
			keys[i] = r.Int63n(maxKey)
		}
		preparedKeys[j] = iter_utils.PartiallyShuffledSlice(keys, shuffleCoef)
	}

	var batchCount atomic.Uint64

	var deleteCount atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(concurrentParallelism)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sl.BulkDelete(preparedKeys[batchCount.Load()%batchSize])
			batchCount.Add(1)

			if deleteCount.Add(1)%compactEvery == 0 {
				sl.Compact()
			}
		}
	})
}
