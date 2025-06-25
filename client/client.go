package client

import (
	"bytes"
	"fmt"
	"net"
	"time"
	"torrent/bitfield"
	"torrent/handshake"
	"torrent/message"
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

func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (*handshake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	req := handshake.New(infoHash, peerID)
	_, err := conn.Write(req.Serialize())
	if err != nil {
		return nil, err
	}

	res, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
		return nil, fmt.Errorf("Expected infoHash %x but got %x", infoHash, res.InfoHash)
	}

	return res, nil
}

func receiveBitfield(conn net.Conn) (bitfield.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}

	if msg == nil {
		err := fmt.Errorf("Expected bitfield but got %v", msg)
		return nil, err
	}

	if msg.ID != message.MsgBitfield { // should be 5 to indicate bitfield msg
		err := fmt.Errorf("Expected ID for bitfield msg but got %v", msg.ID)
		return nil, err
	}

	return msg.Payload, nil
}

func New(peer torrentfile.Peer, infoHash [20]byte, peerID [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	_, err = completeHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := receiveBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{Conn: conn, Choked: true, Bitfield: bf, peer: peer, infoHash: infoHash, peerID: peerID}, nil
}
