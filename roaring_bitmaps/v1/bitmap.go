//go:build goexperiment.simd && amd64

// Package roaringv1 implements Roaring bitmaps with the two classic container
// types: a sorted uint16 array for sparse chunks and a 64Ki-bit bitmap for dense
// ones. It is the baseline of the three versions in this directory; v2 adds run
// containers and v3 adds serialisation on top, so benchmarking the three against
// each other isolates what each addition is worth.
//
// The hot paths run on the shared internal/simdops primitives, which need amd64
// with AVX-512 and GOEXPERIMENT=simd.
package roaringv1

import (
	"slices"
	"strconv"
	"strings"

	"goalgo/roaring_bitmaps/internal/simdops"
)

// Bitmap is a compressed set of uint32 values. The zero value is an empty set
// ready for use. A Bitmap is not safe for concurrent modification.
type Bitmap struct {
	keys  []uint16     // high 16 bits of each chunk, ascending
	conts []*container // chunk contents, parallel to keys
}

// New returns a Bitmap holding the given values.
func New(values ...uint32) *Bitmap {
	b := &Bitmap{}
	b.AddMany(values...)
	return b
}

func hi(v uint32) uint16 { return uint16(v >> 16) }
func lo(v uint32) uint16 { return uint16(v) }

// Add inserts v and reports whether the set changed.
func (b *Bitmap) Add(v uint32) bool {
	k := hi(v)
	i, found := slices.BinarySearch(b.keys, k)
	if !found {
		b.insertAt(i, k, newArrayContainer(4))
	}
	return b.conts[i].add(lo(v))
}

// AddMany inserts every value, reusing the container lookup for runs of values
// that share a key.
func (b *Bitmap) AddMany(values ...uint32) {
	var (
		lastKey  uint16
		lastCont *container
	)
	for _, v := range values {
		k := hi(v)
		if lastCont == nil || k != lastKey {
			i, found := slices.BinarySearch(b.keys, k)
			if !found {
				b.insertAt(i, k, newArrayContainer(4))
			}
			lastKey, lastCont = k, b.conts[i]
		}
		lastCont.add(lo(v))
	}
}

// Remove deletes v and reports whether the set changed.
func (b *Bitmap) Remove(v uint32) bool {
	i, found := slices.BinarySearch(b.keys, hi(v))
	if !found {
		return false
	}

	changed := b.conts[i].remove(lo(v))
	if b.conts[i].card == 0 {
		b.removeAt(i)
	}
	return changed
}

// Contains reports whether v is in the set.
func (b *Bitmap) Contains(v uint32) bool {
	i, found := slices.BinarySearch(b.keys, hi(v))
	return found && b.conts[i].contains(lo(v))
}

// Cardinality returns the number of values in the set.
func (b *Bitmap) Cardinality() int {
	n := 0
	for _, c := range b.conts {
		n += int(c.card)
	}
	return n
}

// IsEmpty reports whether the set holds no values.
func (b *Bitmap) IsEmpty() bool { return len(b.conts) == 0 }

// Clone returns a deep copy of b.
func (b *Bitmap) Clone() *Bitmap {
	out := &Bitmap{
		keys:  slices.Clone(b.keys),
		conts: make([]*container, len(b.conts)),
	}
	for i, c := range b.conts {
		out.conts[i] = c.clone()
	}
	return out
}

// ToArray returns every value in ascending order.
func (b *Bitmap) ToArray() []uint32 {
	out := make([]uint32, 0, b.Cardinality())
	for i, c := range b.conts {
		out = c.appendValues(out, uint32(b.keys[i])<<16)
	}
	return out
}

// And returns the intersection of b and other.
func (b *Bitmap) And(other *Bitmap) *Bitmap {
	out := &Bitmap{}
	i, j := 0, 0
	for i < len(b.keys) && j < len(other.keys) {
		switch {
		case b.keys[i] < other.keys[j]:
			i++
		case b.keys[i] > other.keys[j]:
			j++
		default:
			out.appendContainer(b.keys[i], andContainers(b.conts[i], other.conts[j]))
			i++
			j++
		}
	}
	return out
}

// Or returns the union of b and other.
func (b *Bitmap) Or(other *Bitmap) *Bitmap {
	out := &Bitmap{}
	i, j := 0, 0
	for i < len(b.keys) && j < len(other.keys) {
		switch {
		case b.keys[i] < other.keys[j]:
			out.appendContainer(b.keys[i], b.conts[i].clone())
			i++
		case b.keys[i] > other.keys[j]:
			out.appendContainer(other.keys[j], other.conts[j].clone())
			j++
		default:
			out.appendContainer(b.keys[i], orContainers(b.conts[i], other.conts[j]))
			i++
			j++
		}
	}
	out.appendRest(b, i)
	out.appendRest(other, j)
	return out
}

// AndNot returns the values of b that are not in other.
func (b *Bitmap) AndNot(other *Bitmap) *Bitmap {
	out := &Bitmap{}
	i, j := 0, 0
	for i < len(b.keys) && j < len(other.keys) {
		switch {
		case b.keys[i] < other.keys[j]:
			out.appendContainer(b.keys[i], b.conts[i].clone())
			i++
		case b.keys[i] > other.keys[j]:
			j++
		default:
			out.appendContainer(b.keys[i], andNotContainers(b.conts[i], other.conts[j]))
			i++
			j++
		}
	}
	out.appendRest(b, i)
	return out
}

// Xor returns the symmetric difference of b and other.
func (b *Bitmap) Xor(other *Bitmap) *Bitmap {
	out := &Bitmap{}
	i, j := 0, 0
	for i < len(b.keys) && j < len(other.keys) {
		switch {
		case b.keys[i] < other.keys[j]:
			out.appendContainer(b.keys[i], b.conts[i].clone())
			i++
		case b.keys[i] > other.keys[j]:
			out.appendContainer(other.keys[j], other.conts[j].clone())
			j++
		default:
			out.appendContainer(b.keys[i], xorContainers(b.conts[i], other.conts[j]))
			i++
			j++
		}
	}
	out.appendRest(b, i)
	out.appendRest(other, j)
	return out
}

// Intersects reports whether b and other share any value, without building the
// intersection.
func (b *Bitmap) Intersects(other *Bitmap) bool {
	i, j := 0, 0
	for i < len(b.keys) && j < len(other.keys) {
		switch {
		case b.keys[i] < other.keys[j]:
			i++
		case b.keys[i] > other.keys[j]:
			j++
		default:
			if intersects(b.conts[i], other.conts[j]) {
				return true
			}
			i++
			j++
		}
	}
	return false
}

// AndCardinality returns the size of the intersection without building it.
func (b *Bitmap) AndCardinality(other *Bitmap) int {
	n, i, j := 0, 0, 0
	for i < len(b.keys) && j < len(other.keys) {
		switch {
		case b.keys[i] < other.keys[j]:
			i++
		case b.keys[i] > other.keys[j]:
			j++
		default:
			if a, o := b.conts[i], other.conts[j]; a.kind == kindBitmap && o.kind == kindBitmap {
				n += simdops.AndCardinality(a.bm[:], o.bm[:])
			} else if c := andContainers(a, o); c != nil {
				n += int(c.card)
			}
			i++
			j++
		}
	}
	return n
}

// String renders the set as {1,2,3}, for debugging and test failures.
func (b *Bitmap) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	for i, v := range b.ToArray() {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.FormatUint(uint64(v), 10))
	}
	sb.WriteByte('}')
	return sb.String()
}

func (b *Bitmap) insertAt(i int, k uint16, c *container) {
	b.keys = slices.Insert(b.keys, i, k)
	b.conts = slices.Insert(b.conts, i, c)
}

func (b *Bitmap) removeAt(i int) {
	b.keys = slices.Delete(b.keys, i, i+1)
	b.conts = slices.Delete(b.conts, i, i+1)
}

// appendContainer adds a container under key k, dropping results that came back
// empty. Keys arrive in ascending order, so no search is needed.
func (b *Bitmap) appendContainer(k uint16, c *container) {
	if c == nil || c.card == 0 {
		return
	}
	b.keys = append(b.keys, k)
	b.conts = append(b.conts, c)
}

// appendRest copies the containers of src from index i onward.
func (b *Bitmap) appendRest(src *Bitmap, i int) {
	for ; i < len(src.keys); i++ {
		b.appendContainer(src.keys[i], src.conts[i].clone())
	}
}
