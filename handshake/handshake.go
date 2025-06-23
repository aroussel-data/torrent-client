package handshake

// A Handshake is a special message that a peer uses to identify itself
// see https://www.bittorrent.org/beps/bep_0003.html
type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func New(infoHash [20]byte, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}
