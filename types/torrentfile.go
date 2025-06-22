package types

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (t *TorrentFile) BuildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) RequestPeers(peerID [20]byte, port uint16) ([]Peer, error) {
	url, err := t.BuildTrackerURL(peerID, port)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	btr := BencodeTrackerResponse{}
	err = bencode.Unmarshal(resp.Body, &btr)
	if err != nil {
		return nil, err
	}

	peersList, err := splitPeers([]byte(btr.Peers))
	if err != nil {
		return nil, err
	}

	return peersList, nil
}

func splitPeers(peers []byte) ([]Peer, error) {
	const peerLength = 6 // each peer is 6 bytes (IP + Port)
	buf := []byte(peers)
	if len(buf)%peerLength != 0 {
		err := fmt.Errorf("received malformed list of peers of length %d", len(buf))
		return nil, err
	}
	numPeers := len(buf) / peerLength
	splitPeers := make([]Peer, numPeers)
	for i := range numPeers {
		offset := i * peerLength
		splitPeers[i].IP = net.IP(buf[offset : offset+4])
		splitPeers[i].Port = binary.BigEndian.Uint16(buf[offset+4 : offset+6])
	}

	return splitPeers, nil
}
