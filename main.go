package main

import (
	"fmt"
	"log"
	"os"

	"torrent/client"
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

	// NOTE: just use a random peer ID
	peerID, err := utils.GeneratePeerID()
	utils.FatalCheck(err)

	// NOTE: use the announce URL in torrent file to find peers
	// by calling tracker
	peers, err := torrentFile.RequestPeers(peerID, 6881)
	utils.FatalCheck(err)

	torrent := torrent.Torrent{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    torrentFile.InfoHash,
		PieceHashes: torrentFile.PieceHashes,
		PieceLength: torrentFile.PieceLength,
		Length:      torrentFile.Length,
		Name:        torrentFile.Name,
	}

	client, err := client.New(torrent.Peers[0], torrent.InfoHash, torrent.PeerID)
	utils.FatalCheck(err)

	defer client.Conn.Close()
	log.Printf("Completed handshake with %v", torrent.Peers[0].String())

	// TODO: implement functions to tell peer to unchoke and that we are interested.

}
