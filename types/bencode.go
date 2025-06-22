package types

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

func (info *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *info)
	if err != nil {
		return [20]byte{}, err
	}
	result := sha1.Sum(buf.Bytes())
	return result, nil
}

func (info *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLength := 20
	buf := []byte(info.Pieces)
	if len(buf)%hashLength != 0 {
		err := fmt.Errorf("received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLength
	hashes := make([][20]byte, numHashes)
	for i := range numHashes {
		copy(hashes[i][:], buf[i*hashLength:(i+1)*hashLength])
	}
	return hashes, nil
}

type BencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func (bto BencodeTorrent) ToTorrentFile() (file TorrentFile, err error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	result := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return result, nil
}

type BencodeTrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}
