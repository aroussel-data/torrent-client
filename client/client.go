package client

import (
	"net"
	"torrent/bitfield"
	"torrent/torrentfile"
)

type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield bitfield.Bitfield
	peer     torrentfile.Peer
	infoHash [20]byte
	peerID   [20]byte
}
