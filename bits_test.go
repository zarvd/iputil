package iputil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bitAt(t *testing.T) {
	t.Parallel()

	bytes := []byte{0b1010_1010, 0b0101_0101}
	assert.Equal(t, uint8(1), bitAt(bytes, 0))
	assert.Equal(t, uint8(0), bitAt(bytes, 1))
	assert.Equal(t, uint8(1), bitAt(bytes, 2))
	assert.Equal(t, uint8(0), bitAt(bytes, 3))

	assert.Equal(t, uint8(1), bitAt(bytes, 4))
	assert.Equal(t, uint8(0), bitAt(bytes, 5))
	assert.Equal(t, uint8(1), bitAt(bytes, 6))
	assert.Equal(t, uint8(0), bitAt(bytes, 7))

	assert.Equal(t, uint8(0), bitAt(bytes, 8))
	assert.Equal(t, uint8(1), bitAt(bytes, 9))
	assert.Equal(t, uint8(0), bitAt(bytes, 10))
	assert.Equal(t, uint8(1), bitAt(bytes, 11))

	assert.Equal(t, uint8(0), bitAt(bytes, 12))
	assert.Equal(t, uint8(1), bitAt(bytes, 13))
	assert.Equal(t, uint8(0), bitAt(bytes, 14))
	assert.Equal(t, uint8(1), bitAt(bytes, 15))

	assert.Panics(t, func() { bitAt(bytes, 16) })
}

func Test_setBitAt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name     string
		Initial  []byte
		Index    int
		Bit      uint8
		Expected []byte
	}{
		{
			Name:     "Set first bit to 1",
			Initial:  []byte{0b0000_0000},
			Index:    0,
			Bit:      1,
			Expected: []byte{0b1000_0000},
		},
		{
			Name:     "Set last bit to 1",
			Initial:  []byte{0b0000_0000},
			Index:    7,
			Bit:      1,
			Expected: []byte{0b0000_0001},
		},
		{
			Name:     "Set middle bit to 1",
			Initial:  []byte{0b0000_0000},
			Index:    4,
			Bit:      1,
			Expected: []byte{0b0000_1000},
		},
		{
			Name:     "Set bit to 0",
			Initial:  []byte{0b1111_1111},
			Index:    3,
			Bit:      0,
			Expected: []byte{0b1110_1111},
		},
		{
			Name:     "Set bit in second byte",
			Initial:  []byte{0b0000_0000, 0b0000_0000},
			Index:    9,
			Bit:      1,
			Expected: []byte{0b0000_0000, 0b0100_0000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			setBitAt(tt.Initial, tt.Index, tt.Bit)
			assert.Equal(t, tt.Expected, tt.Initial)
		})
	}

	assert.Panics(t, func() { setBitAt([]byte{0x00}, 8, 1) })
}
