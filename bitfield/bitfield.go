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
