package protoutils

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	
	"github.com/gogo/protobuf/proto"
)

type ProtoMessageReader interface {
	ReadMessage(message proto.Message) error
}

type ProtoReader struct {
	reader *bufio.Reader
}

func NewProtoReader(r io.Reader) *ProtoReader {
	return &ProtoReader{
		reader: bufio.NewReader(r),
	}
}

func (r *ProtoReader) ReadMessage(msg proto.Message) error {
	// 1. read message length (4 bytes, bigEndian)
	var msgLength uint32
	if err := binary.Read(r.reader, binary.BigEndian, &msgLength); err != nil {
		return fmt.Errorf("failed to read message length: %w", err)
	}
	
	// 2. read message body
	msgBytes := make([]byte, msgLength)
	if _, err := io.ReadFull(r.reader, msgBytes); err != nil {
		return fmt.Errorf("failed to read message body: %w", err)
	}
	
	// 3. unmarshal message
	if err := proto.Unmarshal(msgBytes, msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return nil
}
