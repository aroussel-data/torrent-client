package bitfield

// A Bitfield represents the pieces of file that a peer has
type Bitfield []byte

func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex < 0 || byteIndex >= len(bf) {
		return false // index out of bounds
	}
	return bf[byteIndex]>>uint(7-offset)&1 != 0
}

func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex < 0 || byteIndex >= len(bf) {
		return // index out of bounds, do nothing
	}
	bf[byteIndex] |= 1 << uint(7-offset) // set the bit at the offset to 1
}
