package roaringv3

import (
	"encoding/binary"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerializationRoundTrip(t *testing.T) {
	t.Parallel()

	type args struct {
		values      []uint32
		runOptimize bool
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "empty", args: args{values: nil}},
		{name: "single array container", args: args{values: sequence(0, 100)}},
		{name: "single bitmap container", args: args{values: spread(0, 3, 9000)}},
		{name: "run container", args: args{values: sequence(0, 20000), runOptimize: true}},
		{name: "runs below the offset threshold", args: args{values: sequence(0, 20000), runOptimize: true}},
		{
			name: "runs above the offset threshold",
			args: args{
				values:      append(append(sequence(0, 20000), sequence(1<<20, 20000)...), append(sequence(2<<20, 20000), sequence(3<<20, 20000)...)...),
				runOptimize: true,
			},
		},
		{name: "mixed encodings", args: args{values: append(append(sequence(0, 20000), spread(1<<20, 3, 9000)...), sequence(2<<20, 50)...), runOptimize: true}},
		{name: "boundaries of the key space", args: args{values: []uint32{0, 65535, 65536, 1 << 31, ^uint32(0)}}},
		{name: "many chunks", args: args{values: spread(0, 1<<16, 500)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := New(tt.args.values...)
			if tt.args.runOptimize {
				b.RunOptimize()
			}
			want := b.ToArray()

			data := b.ToBytes()
			require.Equal(t, b.SerializedSize(), len(data), "SerializedSize must match ToBytes")

			got, err := FromBytes(data)
			require.NoError(t, err)
			require.Equal(t, want, got.ToArray())
			require.Equal(t, b.Cardinality(), got.Cardinality())

			// A decoded bitmap must serialise back to the identical bytes.
			require.Equal(t, data, got.ToBytes())
		})
	}
}

func TestFromBytesRejectsBadInput(t *testing.T) {
	t.Parallel()

	valid := func() []byte {
		b := New(sequence(0, 20000)...)
		b.RunOptimize()
		return b.ToBytes()
	}

	type args struct {
		data []byte
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "empty input", args: args{data: nil}},
		{name: "truncated cookie", args: args{data: []byte{1, 2}}},
		{name: "unknown cookie", args: args{data: binary.LittleEndian.AppendUint32(nil, 999)}},
		{name: "missing container count", args: args{data: binary.LittleEndian.AppendUint32(nil, 12346)}},
		{name: "truncated body", args: args{data: valid()[:len(valid())-4]}},
		{name: "truncated header", args: args{data: valid()[:6]}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := FromBytes(tt.args.data)

			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}

func TestFromBytesRejectsUnorderedKeys(t *testing.T) {
	t.Parallel()

	b := New(0, 1<<16, 2<<16)
	data := b.ToBytes()

	// The descriptive header starts right after the 8-byte no-run preamble; swap
	// the first two keys so they no longer ascend.
	first := binary.LittleEndian.Uint16(data[8:])
	second := binary.LittleEndian.Uint16(data[12:])
	binary.LittleEndian.PutUint16(data[8:], second)
	binary.LittleEndian.PutUint16(data[12:], first)

	got, err := FromBytes(data)

	require.ErrorContains(t, err, "out of order")
	require.Nil(t, got)
}

func TestRunOptimizeShrinksSerializedSize(t *testing.T) {
	t.Parallel()

	b := New(sequence(0, 40000)...)
	before := b.SerializedSize()

	require.True(t, b.RunOptimize())

	require.Less(t, b.SerializedSize(), before)
}

func FuzzSerializationRoundTrip(f *testing.F) {
	f.Add(uint64(1), uint32(500), false)
	f.Add(uint64(3), uint32(70000), true)

	f.Fuzz(func(t *testing.T, seed uint64, span uint32, runOptimize bool) {
		rng := rand.New(rand.NewPCG(seed, seed+1))
		b := New(randomRuns(rng, int(span%500), span|1)...)
		if runOptimize {
			b.RunOptimize()
		}
		want := b.ToArray()

		data := b.ToBytes()
		got, err := FromBytes(data)

		require.NoError(t, err)
		require.Equal(t, want, got.ToArray())
		require.Equal(t, data, got.ToBytes())
	})
}

func FuzzFromBytesDoesNotPanic(f *testing.F) {
	f.Add(New(sequence(0, 100)...).ToBytes())
	f.Add([]byte{0, 0, 0, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		b, err := FromBytes(data)
		if err == nil {
			// Whatever parsed must at least be self-consistent.
			require.Equal(t, len(b.ToArray()), b.Cardinality())
		}
	})
}
