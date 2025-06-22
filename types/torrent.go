package types

type Torrent struct {
	Peers       []Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

// TODO: add the download() method
func (t *Torrent) Download() ([]byte, error) {
	return nil, nil
}
