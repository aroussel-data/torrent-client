package main

import (
	"fmt"
	"os"

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
	fmt.Println(torrentFile.RequestPeers(peerID, 6881))

}
