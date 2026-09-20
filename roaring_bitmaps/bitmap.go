// Package roaring_bitmaps is a Roaring bitmap: a compressed set of uint32 values
// split into 65536-value chunks, each stored in whichever of three encodings is
// cheapest for its contents — a sorted uint16 array, a 64Ki-bit bitmap, or a
// list of runs.
//
// Run encoding is never chosen on its own. Containers keep the shape an
// operation produced until RunOptimize is called, which is the same bargain the
// reference C implementation strikes: scanning for runs costs a pass over the
// data, so the caller decides when it is worth paying. Operations between two
// run containers stay in interval arithmetic and never materialise the values.
//
// ToBytes/FromBytes speak the portable format shared with the C and Java
// implementations; Rank/Select answer positional queries from the containers
// plus a lazily built prefix-sum cache.
//
// The hot paths run on internal/simdops, which is AVX-512 under
// GOEXPERIMENT=simd on amd64 and plain Go otherwise.
package roaring_bitmaps

import (
	"slices"
	"strconv"
	"strings"
	"unsafe"

	"goalgo/roaring_bitmaps/internal/simdops"
)

// Bitmap is a compressed set of uint32 values. The zero value is an empty set
// ready for use. A Bitmap is not safe for concurrent modification.
type Bitmap struct {
	keys  []uint16     // high 16 bits of each chunk, ascending
	conts []*container // chunk contents, parallel to keys

	// prefix[i] is the number of values in the first i containers, built on
	// demand by Rank/Select and dropped by any change to the contents. It turns
	// both queries from a walk over the containers into a binary search.
	prefix []int
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
	b.prefix = nil
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
	b.prefix = nil
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
	b.prefix = nil
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

// RunOptimize re-encodes every chunk as runs where that is the smallest of the
// three representations, and reports whether any chunk changed. It is the only
// path that produces run containers.
func (b *Bitmap) RunOptimize() bool {
	changed := false
	for _, c := range b.conts {
		before := c.kind
		c.runOptimize()
		changed = changed || c.kind != before
	}
	return changed
}

// Stats reports how many chunks use each encoding, for benchmarks and for seeing
// what RunOptimize actually did.
func (b *Bitmap) Stats() (arrays, bitmaps, runs int) {
	for _, c := range b.conts {
		switch c.kind {
		case kindArray:
			arrays++
		case kindBitmap:
			bitmaps++
		case kindRun:
			runs++
		}
	}
	return arrays, bitmaps, runs
}

// SizeInBytes reports the heap memory the bitmap occupies: the key and
// container slices, every container struct with its payload, and the
// Rank/Select prefix cache once one has been built. Slice capacities count
// rather than lengths, since capacity is what was allocated.
func (b *Bitmap) SizeInBytes() int {
	n := int(unsafe.Sizeof(*b)) +
		cap(b.keys)*int(unsafe.Sizeof(uint16(0))) +
		cap(b.conts)*int(unsafe.Sizeof((*container)(nil))) +
		cap(b.prefix)*int(unsafe.Sizeof(int(0)))
	for _, c := range b.conts {
		n += int(unsafe.Sizeof(*c)) +
			cap(c.arr)*int(unsafe.Sizeof(uint16(0))) +
			cap(c.runs)*int(unsafe.Sizeof(interval{}))
		if c.bm != nil {
			n += bitmapBytes
		}
	}
	return n
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
