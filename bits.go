package iputil

func indexOfBits(i int) int {
	const BitSize = 8
	return (BitSize - 1) - i%BitSize
}

// setBitAt sets the bit at the i-th position in the byte slice to the given value.
// Panics if the index is out of bounds.
// For example,
// - setBitAt([0x00, 0x00], 8, 1) returns [0x00, 0b1000_0000].
// - setBitAt([0xff, 0xff], 0, 0) returns [0b0111_1111, 0xff].
func setBitAt(bytes []byte, i int, bit uint8) {
	if bit == 1 {
		bytes[i/8] |= 1 << indexOfBits(i)
	} else {
		bytes[i/8] &^= 1 << indexOfBits(i)
	}
}

// bitAt returns the bit at the i-th position in the byte slice.
// The return value is either 0 or 1 as uint8.
// Panics if the index is out of bounds.
func bitAt(bytes []byte, i int) uint8 {
	return bytes[i/8] >> indexOfBits(i) & 1
}
