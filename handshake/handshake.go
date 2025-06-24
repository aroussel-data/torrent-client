package handshake

import (
	"fmt"
	"io"
)

// A Handshake is a special message that a peer uses to identify itself
// see https://www.bittorrent.org/beps/bep_0003.html
type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

// Serialize the handshake into a format that can be sent to peer
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], h.Pstr)
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

// Read the handshake response from a connection
func Read(r io.Reader) (*Handshake, error) {
	lengthbuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthbuf)
	if err != nil {
		return nil, err
	}
	pstrlen := int(lengthbuf[0])

	if pstrlen == 0 {
		err := fmt.Errorf("invalid protocol identifier length. should be 19")
		return nil, err
	}

	// buffer is going to be 48 + pstrlen bc we already read first
	// byte to get lengthbuf
	handshakebuf := make([]byte, 48+pstrlen)
	_, err = io.ReadFull(r, handshakebuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	// copy from 19+8 bytes bc there are 8 zeroed bytes after protocol identifier
	copy(infoHash[:], handshakebuf[pstrlen+8:pstrlen+8+20]) // infoHash is next 20 bytes
	copy(peerID[:], handshakebuf[pstrlen+8+20:])            // rest of buffer will be the peerID

	return &Handshake{
		Pstr:     string(handshakebuf[:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}, nil
}

func New(infoHash [20]byte, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}
