package roaring_bitmaps

import (
	"encoding/binary"
	"fmt"

	"goalgo/roaring_bitmaps/internal/simdops"
)

// The portable Roaring format, as shared with the C and Java implementations.
//
// A stream starts with a cookie. When no chunk uses run encoding the cookie is
// serialCookieNoRun followed by a uint32 container count; otherwise the count is
// packed into the cookie word and a bitset of run flags follows, one bit per
// container. Then comes a descriptive header of (key, cardinality-1) pairs, an
// optional table of absolute byte offsets, and finally the container bodies.
// Everything is little-endian.
const (
	serialCookie      = 12347
	serialCookieNoRun = 12346

	// noOffsetThreshold is the container count below which a run-flagged stream
	// may omit the offset table.
	noOffsetThreshold = 4
)

// ToBytes serialises the bitmap in the portable Roaring format.
func (b *Bitmap) ToBytes() []byte {
	n := len(b.conts)
	hasRun := false
	for _, c := range b.conts {
		if c.kind == kindRun {
			hasRun = true
			break
		}
	}

	buf := make([]byte, 0, b.serializedSize(hasRun))

	if hasRun {
		buf = binary.LittleEndian.AppendUint32(buf, uint32(serialCookie|((n-1)<<16)))
		flags := make([]byte, (n+7)/8)
		for i, c := range b.conts {
			if c.kind == kindRun {
				flags[i/8] |= 1 << (i % 8)
			}
		}
		buf = append(buf, flags...)
	} else {
		buf = binary.LittleEndian.AppendUint32(buf, serialCookieNoRun)
		buf = binary.LittleEndian.AppendUint32(buf, uint32(n))
	}

	for i, c := range b.conts {
		buf = binary.LittleEndian.AppendUint16(buf, b.keys[i])
		buf = binary.LittleEndian.AppendUint16(buf, uint16(c.card-1))
	}

	if withOffsets := !hasRun || n >= noOffsetThreshold; withOffsets {
		offset := len(buf) + 4*n
		for _, c := range b.conts {
			buf = binary.LittleEndian.AppendUint32(buf, uint32(offset))
			offset += c.serializedSize()
		}
	}

	for _, c := range b.conts {
		buf = c.appendTo(buf)
	}
	return buf
}

// SerializedSize returns the number of bytes ToBytes will produce.
func (b *Bitmap) SerializedSize() int {
	for _, c := range b.conts {
		if c.kind == kindRun {
			return b.serializedSize(true)
		}
	}
	return b.serializedSize(false)
}

func (b *Bitmap) serializedSize(hasRun bool) int {
	n := len(b.conts)

	size := 4 // cookie
	if hasRun {
		size += (n + 7) / 8
	} else {
		size += 4 // container count
	}
	size += 4 * n // descriptive header
	if !hasRun || n >= noOffsetThreshold {
		size += 4 * n // offset header
	}
	for _, c := range b.conts {
		size += c.serializedSize()
	}
	return size
}

// FromBytes parses a bitmap in the portable Roaring format, running on
// internal/simdops.
func FromBytes(data []byte) (*Bitmap, error) {
	return FromBytesWithOps(defaultOps, data)
}

// FromBytesWithOps parses a bitmap in the portable Roaring format that runs
// every container operation through ops.
func FromBytesWithOps(ops Ops, data []byte) (*Bitmap, error) {
	r := &reader{data: data}

	cookieWord, err := r.uint32()
	if err != nil {
		return nil, fmt.Errorf("roaring: reading cookie: %w", err)
	}

	var (
		n      int
		hasRun bool
	)
	switch cookie := int(cookieWord & 0xFFFF); cookie {
	case serialCookieNoRun:
		count, err := r.uint32()
		if err != nil {
			return nil, fmt.Errorf("roaring: reading container count: %w", err)
		}
		n = int(count)
	case serialCookie:
		hasRun = true
		n = int(cookieWord>>16) + 1
	default:
		return nil, fmt.Errorf("roaring: unknown cookie %d", cookie)
	}

	if n < 0 || n > 1<<16 {
		return nil, fmt.Errorf("roaring: implausible container count %d", n)
	}

	runFlags := make([]byte, 0)
	if hasRun {
		if runFlags, err = r.bytes((n + 7) / 8); err != nil {
			return nil, fmt.Errorf("roaring: reading run flags: %w", err)
		}
	}

	b := &Bitmap{
		keys:  make([]uint16, n),
		conts: make([]*container, n),
		simd:  ops,
	}
	cards := make([]int, n)
	for i := 0; i < n; i++ {
		key, err := r.uint16()
		if err != nil {
			return nil, fmt.Errorf("roaring: reading key %d: %w", i, err)
		}
		card, err := r.uint16()
		if err != nil {
			return nil, fmt.Errorf("roaring: reading cardinality %d: %w", i, err)
		}
		b.keys[i] = key
		cards[i] = int(card) + 1
	}

	if !hasRun || n >= noOffsetThreshold {
		if _, err := r.bytes(4 * n); err != nil {
			return nil, fmt.Errorf("roaring: reading offset header: %w", err)
		}
	}

	for i := 0; i < n; i++ {
		isRun := hasRun && runFlags[i/8]&(1<<(i%8)) != 0
		c, err := readContainer(ops, r, cards[i], isRun)
		if err != nil {
			return nil, fmt.Errorf("roaring: reading container %d: %w", i, err)
		}
		b.conts[i] = c
	}

	if err := b.validate(); err != nil {
		return nil, err
	}
	return b, nil
}

// validate rejects streams whose keys are not strictly ascending, which every
// other operation relies on.
func (b *Bitmap) validate() error {
	for i := 1; i < len(b.keys); i++ {
		if b.keys[i] <= b.keys[i-1] {
			return fmt.Errorf("roaring: container keys out of order at %d", i)
		}
	}
	return nil
}

func (c *container) serializedSize() int {
	switch c.kind {
	case kindArray:
		return 2 * int(c.card)
	case kindRun:
		return 2 + 4*len(c.runs)
	default:
		return simdops.BitmapWords * 8
	}
}

func (c *container) appendTo(buf []byte) []byte {
	switch c.kind {
	case kindArray:
		for _, v := range c.arr {
			buf = binary.LittleEndian.AppendUint16(buf, v)
		}
	case kindRun:
		buf = binary.LittleEndian.AppendUint16(buf, uint16(len(c.runs)))
		for _, iv := range c.runs {
			buf = binary.LittleEndian.AppendUint16(buf, iv.start)
			// The format stores length-1, so a one-value run writes a zero.
			buf = binary.LittleEndian.AppendUint16(buf, iv.last-iv.start)
		}
	default:
		for _, w := range c.bm[:] {
			buf = binary.LittleEndian.AppendUint64(buf, w)
		}
	}
	return buf
}

func readContainer(ops Ops, r *reader, card int, isRun bool) (*container, error) {
	switch {
	case isRun:
		nruns, err := r.uint16()
		if err != nil {
			return nil, err
		}
		runs := make([]interval, nruns)
		for i := range runs {
			start, err := r.uint16()
			if err != nil {
				return nil, err
			}
			length, err := r.uint16()
			if err != nil {
				return nil, err
			}
			if int(start)+int(length) > 0xFFFF {
				return nil, fmt.Errorf("run %d overflows the container", i)
			}
			runs[i] = interval{start: start, last: start + length}
			if i > 0 && int(runs[i-1].last)+1 >= int(start) {
				return nil, fmt.Errorf("run %d overlaps or touches its predecessor", i)
			}
		}
		c := newRunContainer(runs)
		if int(c.card) != card {
			return nil, fmt.Errorf("run cardinality %d does not match header %d", c.card, card)
		}
		return c, nil

	case card > arrayMax:
		words, err := r.bytes(simdops.BitmapWords * 8)
		if err != nil {
			return nil, err
		}
		c := newBitmapContainer()
		for i := range c.bm[:] {
			c.bm[i] = binary.LittleEndian.Uint64(words[8*i:])
		}
		c.card = int32(ops.Popcount(c.bm[:]))
		if int(c.card) != card {
			return nil, fmt.Errorf("bitmap cardinality %d does not match header %d", c.card, card)
		}
		return c, nil

	default:
		c := newArrayContainer(card)
		c.arr = c.arr[:card]
		for i := 0; i < card; i++ {
			v, err := r.uint16()
			if err != nil {
				return nil, err
			}
			if i > 0 && v <= c.arr[i-1] {
				return nil, fmt.Errorf("array values out of order at %d", i)
			}
			c.arr[i] = v
		}
		c.card = int32(card)
		return c, nil
	}
}

// reader walks a byte slice, reporting a single error kind for any short read.
type reader struct {
	data []byte
	pos  int
}

var errShortStream = fmt.Errorf("unexpected end of stream")

func (r *reader) bytes(n int) ([]byte, error) {
	if r.pos+n > len(r.data) {
		return nil, errShortStream
	}
	out := r.data[r.pos : r.pos+n]
	r.pos += n
	return out, nil
}

func (r *reader) uint16() (uint16, error) {
	b, err := r.bytes(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func (r *reader) uint32() (uint32, error) {
	b, err := r.bytes(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}
