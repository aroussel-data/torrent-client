package client

import (
	"net"
	"torrent/bitfield"
	"torrent/handshake"
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

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (handshake.Handshake, error) {
	return handshake.Handshake{}, nil
}

func receiveBitfield(conn net.Conn) (bitfield.Bitfield, error) {
	return nil, nil
}

func New(peer torrentfile.Peer, infoHash [20]byte, peerID [20]byte) *Client {
	return &Client{peer: peer, infoHash: infoHash, peerID: peerID}
}
