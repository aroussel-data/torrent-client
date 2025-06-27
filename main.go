package main

import (
	"fmt"
	"os"

	"torrent/torrent"
	"torrent/utils"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "expected a torrent file")
		os.Exit(1)
	}
	path := args[1]

	torrentFile, err := utils.Open(path)
	utils.FatalCheck(err)

	// just use a random peer ID
	peerID, err := utils.GeneratePeerID()
	utils.FatalCheck(err)

	// use the announce URL in torrent file to find peers
	// by calling tracker
	peers, err := torrentFile.RequestPeers(peerID, 6881)
	utils.FatalCheck(err)

	tr := torrent.Torrent{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    torrentFile.InfoHash,
		PieceHashes: torrentFile.PieceHashes,
		PieceLength: torrentFile.PieceLength,
		Length:      torrentFile.Length,
		Name:        torrentFile.Name,
	}
	// fmt.Printf("Torrent has %d piece hashes and each has %d piece length\n", len(tr.PieceHashes), tr.PieceLength)
	// that makes a the total 700MB file

	res, err := tr.Download()
	utils.FatalCheck(err)
	fmt.Printf("Downloaded %d bytes from torrent %s\n", len(res), tr.Name)

	// TODO: allow user to specify output file name
	file, err := os.Create(tr.Name)
	utils.FatalCheck(err)
	defer file.Close()
	_, err = file.Write(res)
	utils.FatalCheck(err)
}
