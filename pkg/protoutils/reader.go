package protoutils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
)

type ProtoMessageReader interface {
	io.Closer
	ReadMessage(message proto.Message) error
}

type ProtoReader struct {
	reader io.ReadCloser
}

func NewProtoReader(r io.ReadCloser) *ProtoReader {
	return &ProtoReader{
		reader: r,
	}
}

func (r *ProtoReader) ReadMessage(msg proto.Message) error {
	// 1. read message length (4 bytes, bigEndian)
	var msgLength uint32
	if err := binary.Read(r.reader, binary.BigEndian, &msgLength); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to read message length: %w", err)
	}

	// 2. read message body
	msgBytes := make([]byte, msgLength)
	if _, err := io.ReadFull(r.reader, msgBytes); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to read message body: %w", err)
	}

	// 3. unmarshal message
	if err := proto.Unmarshal(msgBytes, msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return nil
}

func (r *ProtoReader) Close() error {
	return r.reader.Close()
}
