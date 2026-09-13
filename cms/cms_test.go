package cms

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"
)

// Вставляет 1M элементов (1000 уникальных) в CMS с шириной 1024 и глубиной 10, используя 3 хэш-функции.

func repeat(n, k int) []uint64 {
	counts := make([]uint64, n)
	for i := 0; i < n; i++ {
		counts[i] = uint64(k)
	}
	return counts
}

func TestCMS_1024_1M_100U(t *testing.T) {
	hash1 := MakeHash(1)
	hash2 := MakeHash(2)
	hash3 := MakeHash(3)
	cms := NewCMS(1024, []func([]byte) uint64{hash1, hash2, hash3})
	data := make([][]byte, 1_000_000)
	M := 1000
	for k := 0; k < 1000; k++ {
		for i := 0; i < M; i++ {
			data[k*M+i] = []byte(fmt.Sprintf("word%d", i))
			cms.Add(data[k*M+i])
		}
	}
	actualCounts := repeat(1000, 1000)
	md := MaxDeviation(cms, data, actualCounts)
	fmt.Println("Precision:", cms.Precision(), "Max Deviation:", md)
}

func readFileByWords(filename string) ([]string, error) {
	content, err := fs.ReadFile(fs.FS(os.DirFS(".")), filename)
	if err != nil {
		return nil, err
	}
	words := strings.Fields(string(content))
	return words, nil
}

func BenchmarkCMS_1024_WAP(b *testing.B) {
	hash1 := MakeHash(1)
	hash2 := MakeHash(2)
	hash3 := MakeHash(3)
	cms := NewCMS(1024, []func([]byte) uint64{hash1, hash2, hash3})
	// Load War and Peace words into CMS
	words, err := readFileByWords("2600-0_lemmatized.txt")
	if err != nil {
		b.Fatal(err)
	}
	for _, word := range words {
		cms.Add([]byte(word))
	}

	b.Run("WAP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, word := range words {
				cms.CountMin([]byte(word))
			}
		}
	})
}

// count of war and peace words

func BenchmarkMapCount_WAP(b *testing.B) {
	wordCounts := make(map[string]uint64)
	words, err := readFileByWords("2600-0_lemmatized.txt")
	if err != nil {
		b.Fatal(err)
	}
	b.Run("WAP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, word := range words {
				wordCounts[word]++
			}
		}
	})
}
