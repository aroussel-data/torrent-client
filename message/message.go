package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
	// MsgChoke chokes the receiver
	MsgChoke messageID = 0
	// MsgUnchoke unchokes the receiver
	MsgUnchoke messageID = 1
	// MsgInterested expresses interest in receiving data
	MsgInterested messageID = 2
	// MsgNotInterested expresses disinterest in receiving data
	MsgNotInterested messageID = 3
	// MsgHave alerts the receiver that the sender has downloaded a piece
	MsgHave messageID = 4
	// MsgBitfield encodes which pieces that the sender has downloaded
	MsgBitfield messageID = 5
	// MsgRequest requests a block of data from the receiver
	MsgRequest messageID = 6
	// MsgPiece delivers a block of data to fulfill a request
	MsgPiece messageID = 7
	// MsgCancel cancels a request
	MsgCancel messageID = 8
)

// Message stores ID and payload of a message
type Message struct {
	ID      messageID
	Payload []byte
}

// Serializes a msg into a buffer of:
// [<length><msgID><payload>]
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1)         // 1 is extra byte for message ID
	buf := make([]byte, length+4)                // 4 extra bytes for the length prefix (all ints sent in protocol are four bytes)
	binary.BigEndian.PutUint32(buf[0:4], length) // put length prefix at start of buffer
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

func FormatRequest(index int, offset int, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))  // 4 bytes for index
	binary.BigEndian.PutUint32(payload[4:8], uint32(offset)) // 4 bytes for offset
	binary.BigEndian.PutUint32(payload[8:], uint32(length))  // 4 bytes for length
	// binary.BigEndian.PutUint32(payload[8:12], uint32(length)) // 4 bytes for length
	return &Message{ID: MsgRequest, Payload: payload}
}

func ParseHave(msg *Message) (int, error) {
	if msg.ID != MsgHave || len(msg.Payload) < 4 {
		return -1, fmt.Errorf("invalid message ID or not enough payload")
	}
	index := binary.BigEndian.Uint32(msg.Payload)
	return int(index), nil
}

func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != MsgPiece || len(msg.Payload) < 8 {
		return -1, fmt.Errorf("invalid message ID or not enough payload")
	}
	// first 4 bytes are the index of the piece, next 4 bytes are the offset
	pieceIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if pieceIndex != index {
		return -1, fmt.Errorf("Expected piece index %d, got %d", index, pieceIndex)
	}
	offset := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if offset >= len(buf) {
		return -1, fmt.Errorf("offset %d is out of bounds for buffer of length %d", offset, len(buf))
	}
	data := msg.Payload[8:] // the actual data starts after the first 8 bytes
	if len(data)+offset > len(buf) {
		return -1, fmt.Errorf("data length %d with offset %d exceeds buffer length %d", len(data), offset, len(buf))
	}
	// copy the payload into the buffer at the offset
	copy(buf[offset:], data)
	return len(data), nil
}

// Read parses a message from a stream. Returns `nil` on keep-alive message
func Read(r io.Reader) (*Message, error) {
	// from docs: peer wire protocol consists of a handshake followed by a never-ending stream of length-prefixed messages
	// All later integers sent in the protocol are encoded as four bytes big-endian.
	// That's it for handshaking, next comes an alternating stream of length prefixes and messages. Messages of length zero are keepalives, and ignored. Keepalives are generally sent once every two minutes, but note that timeouts can be done much more quickly when data is expected.
	lengthBuf := make([]byte, 4) // this is the length prefix for all messages
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf) // length indicator for the actual msg is 4 byte buffer

	// keep-alive message
	if length == 0 {
		return nil, nil
	}

	// read actual message
	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}
