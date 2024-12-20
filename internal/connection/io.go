package connection

import (
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protoutils"
	"io"
)

type StreamReadWriter struct {
	rawRwc io.ReadWriteCloser
	pr     *protoutils.ProtoReader
	pw     *protoutils.ProtoWriter
}

func NewStreamReadWriter(rwc io.ReadWriteCloser) *StreamReadWriter {
	pr := protoutils.NewProtoReader(rwc)
	pw := protoutils.NewProtoWriter(rwc)
	return &StreamReadWriter{
		rawRwc: rwc,
		pr:     pr,
		pw:     pw,
	}
}

func (rw *StreamReadWriter) Read(p []byte) (n int, err error) {
	return rw.rawRwc.Read(p)
}

func (rw *StreamReadWriter) Write(p []byte) (n int, err error) {
	return rw.rawRwc.Write(p)
}

func (rw *StreamReadWriter) ReadMessage(message proto.Message) error {
	err := rw.pr.ReadMessage(message)
	if err != nil {
		return errorutils.ErrStreamReadMessage(err)
	}

	return nil
}

func (rw *StreamReadWriter) WriteMessage(m proto.Message) error {
	err := rw.pw.WriteMessage(m)
	if err != nil {
		return errorutils.ErrStreamWriteMessage(err)
	}

	return nil
}

func (rw *StreamReadWriter) Close() error {
	return rw.rawRwc.Close()
}
