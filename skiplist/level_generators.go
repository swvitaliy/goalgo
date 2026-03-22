package skiplist

import (
	"time"

	"golang.org/x/exp/rand"
)

// Генераторы уровней для вставки ключей в skiplist

// RandomGenerator - рандомный генератор для случайной вставки
var RandomGenerator = NewRandomLevelGenerator(uint64(time.Now().UnixNano()))

// BatchFriendlyGenerator - генератор для вставки ключей, которые идут в порядке возрастания (или близко к нему).
// Такой генератор будет выдавать уровни, которые не слишком сильно отличаются от предыдущего,
// что поможет избежать длинных "хвостов" в skiplist при вставке отсортированных данных.
var BatchFriendlyGenerator = NewBatchLevelGenerator(uint64(time.Now().UnixNano()))

// RandomLevelGenerator - генератор уровней для случайной вставки
type RandomLevelGenerator struct {
	randSeed *rand.Rand
}

// NewRandomLevelGenerator создает новый генератор с seed
func NewRandomLevelGenerator(seed uint64) *RandomLevelGenerator {
	return &RandomLevelGenerator{
		randSeed: rand.New(rand.NewSource(seed)),
	}
}

// NextLevel возвращает следующий уровень для вставляемого ключа
func (_ *RandomLevelGenerator) NextLevel() uint64 {
	var lvl uint64 = 1
	for rand.Float64() < probability && lvl < maxLevel {
		lvl++
	}
	return lvl
}

// BatchLevelGenerator генерирует уровни для bulk вставки
type BatchLevelGenerator struct {
	randSeed  *rand.Rand
	prevLevel uint64
}

// NewBatchLevelGenerator создает новый генератор с seed
func NewBatchLevelGenerator(seed uint64) *BatchLevelGenerator {
	return &BatchLevelGenerator{
		randSeed:  rand.New(rand.NewSource(seed)),
		prevLevel: 1,
	}
}

// NextLevel возвращает следующий уровень для вставляемого ключа
func (b *BatchLevelGenerator) NextLevel() uint64 {
	// Основная идея: небольшие колебания вокруг предыдущего уровня
	// чтобы соседние ключи имели похожие высоты
	var level = b.prevLevel

	// Иногда немного увеличиваем уровень
	if b.randSeed.Float64() < probability && level < maxLevel {
		level++
	}

	// Иногда немного уменьшаем уровень
	if b.randSeed.Float64() < (1.0-probability) && level > 1 {
		level--
	}

	b.prevLevel = level
	return level
}
