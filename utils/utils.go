package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"torrent/types"

	"github.com/jackpal/bencode-go"
)

func Check(er error) {
	if er != nil {
		fmt.Fprintln(os.Stderr, er)
	}
}

func FatalCheck(er error) {
	if er != nil {
		log.Fatal(er)
	}
}

func Open(path string) (types.TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return types.TorrentFile{}, err
	}
	defer file.Close()

	bto := types.BencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return types.TorrentFile{}, err
	}
	return bto.ToTorrentFile()
}

func GeneratePeerID() ([20]byte, error) {
	var peerId [20]byte
	_, err := rand.Read(peerId[:])
	if err != nil {
		return peerId, err
	}
	return peerId, nil
}
