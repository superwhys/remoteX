package protoutils

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
)

type ProtoMessageWriter interface {
	io.Closer
	WriteMessage(m proto.Message) error
}

type ProtoWriter struct {
	writer io.WriteCloser
}

func NewProtoWriter(w io.WriteCloser) *ProtoWriter {
	return &ProtoWriter{
		writer: w,
	}
}

func (w *ProtoWriter) WriteMessage(msg proto.Message) error {
	// 1. marshal proto message
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 2. calc message length
	msgLength := uint32(len(msgBytes))
	// 3. send length data (4 bytes, bigEndian)
	if err := binary.Write(w.writer, binary.BigEndian, msgLength); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	// 4. send proto message
	if _, err := w.writer.Write(msgBytes); err != nil {
		return fmt.Errorf("failed to write message body: %w", err)
	}

	return nil
}

func (w *ProtoWriter) Close() error {
	return w.writer.Close()
}
