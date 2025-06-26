package torrent

import (
	"fmt"
	"torrent/client"
	"torrent/torrentfile"
)

type Torrent struct {
	Peers       []torrentfile.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

func (t *Torrent) calculatePieceLength(index int) int {
	// if the end of the piece is beyond the total length of the torrent
	// then the piece length is the remaining length of the torrent
	begin := t.PieceLength * index
	end := min(begin+t.PieceLength, t.Length)
	return end - begin
}

func (t *Torrent) startDownloadFromPeer(peer torrentfile.Peer, workQueue chan *pieceWork, results chan *pieceResult) {
	client, err := client.New(peer, t.InfoHash, t.PeerID)
	if err != nil {
		fmt.Printf("Failed to connect to peer %s: %v\n", peer.String(), err)
		return // handle error appropriately
	}
	defer client.Conn.Close()

	client.SendUnchoke()
	client.SendInterested()

	for work := range workQueue {
		if !client.Bitfield.HasPiece(work.index) {
			fmt.Printf("Peer %s does not have piece at index %d\n", peer.String(), work.index)
			workQueue <- work // requeue the work if peer does not have the piece
			continue
		}
		buf, err := client.RequestPiece(work.index, work.length)
		if err != nil {
			fmt.Printf("Failed to download piece at index %d from peer %s: %v\n", work.index, peer.String(), err)
			workQueue <- work // requeue the work if there was an error
			continue          // handle error appropriately
		}

		results <- &pieceResult{index: work.index, buf: buf}
	}
}

func (t *Torrent) Download() ([]byte, error) {
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)

	for i, hash := range t.PieceHashes {
		length := t.calculatePieceLength(i) // need to do this in case we are at last piece
		workQueue <- &pieceWork{
			index:  i,
			hash:   hash,
			length: length,
		}
	}

	// for each peer we want to attempt to download a piece from the queue if peer has it
	for _, peer := range t.Peers {
		go t.startDownloadFromPeer(peer, workQueue, results)
	}

	return nil, nil
}
