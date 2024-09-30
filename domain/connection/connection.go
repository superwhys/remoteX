package connection

import (
	"crypto/tls"
	"io"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/counter"
	"github.com/superwhys/remoteX/pkg/limiter"
	"github.com/superwhys/remoteX/pkg/protocol"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type BaseConnection interface {
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}

type Stream interface {
	io.ReadWriteCloser
	BaseConnection

	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter

	GetNodeId() common.NodeID
	SetNodeId(nodeId common.NodeID)
}

type StreamConnection interface {
	io.Closer
	BaseConnection

	AcceptStream() (Stream, error)
	OpenStream() (Stream, error)

	ConnectionState() tls.ConnectionState
}

type TlsConn interface {
	StreamConnection

	IsServer() bool
	GetConnectionId() string
	SetStatus(status protocol.ConnectionStatus)
	UpdateLastHeartbeat()

	GetNodeId() common.NodeID
	SetNodeId(nodeId common.NodeID)
}

var _ TlsConn = (*InternalConn)(nil)

type InternalConn struct {
	StreamConnection
	*Connection
	isServer bool
}

func NewInternalConnection(sc StreamConnection, conn *Connection, isServer bool) *InternalConn {
	return &InternalConn{sc, conn, isServer}
}

func (i *InternalConn) IsServer() bool {
	return i.isServer
}

func (i *InternalConn) GetNodeId() common.NodeID {
	return i.Connection.NodeId
}

func (i *InternalConn) AcceptStream() (Stream, error) {
	stream, err := i.StreamConnection.AcceptStream()
	if err != nil {
		return nil, err
	}

	stream.SetNodeId(i.GetNodeId())
	return stream, nil
}

func (i *InternalConn) OpenStream() (Stream, error) {
	stream, err := i.StreamConnection.OpenStream()
	if err != nil {
		return nil, err
	}

	stream.SetNodeId(i.GetNodeId())
	return stream, nil
}

func (c *Connection) SetStatus(status protocol.ConnectionStatus) {
	c.Status = status
}

func (c *Connection) UpdateLastHeartbeat() {
	c.LastHeartbeat = time.Now().Unix()
}

func (c *Connection) SetNodeId(nodeId common.NodeID) {
	c.NodeId = nodeId
}

var _ Stream = (*LimiterStream)(nil)

type LimiterStream struct {
	Stream
	rd *limiter.LimitReader
	wr *limiter.LimitWriter
	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter
}

func PackLimiterStream(stream Stream, limiter *limiter.Limiter) *LimiterStream {
	rd, wr := limiter.GetNodeRateLimiter(stream)
	return &LimiterStream{
		Stream:             stream,
		rd:                 rd,
		wr:                 wr,
		ProtoMessageReader: protoutils.NewProtoReader(rd),
		ProtoMessageWriter: protoutils.NewProtoWriter(wr),
	}
}

// Read rewrite the method to use LimiterReader
func (l *LimiterStream) Read(p []byte) (n int, err error) {
	return l.rd.Read(p)
}

// Write rewrite the method to use LimiterWriter
func (l *LimiterStream) Write(p []byte) (n int, err error) {
	return l.wr.Write(p)
}

func (l *LimiterStream) ReadMessage(message proto.Message) error {
	return l.ProtoMessageReader.ReadMessage(message)
}

func (l *LimiterStream) WriteMessage(m proto.Message) error {
	return l.ProtoMessageWriter.WriteMessage(m)
}

var _ Stream = (*CounterStream)(nil)

type CounterStream struct {
	Stream
	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter

	rd *counter.CountingReader
	wr *counter.CountingWriter
}

func PackCounterConnection(stream Stream, rd io.Reader, wr io.Writer) *CounterStream {
	crd := &counter.CountingReader{Reader: rd}
	cwr := &counter.CountingWriter{Writer: wr}

	return &CounterStream{
		Stream:             stream,
		rd:                 crd,
		wr:                 cwr,
		ProtoMessageReader: protoutils.NewProtoReader(crd),
		ProtoMessageWriter: protoutils.NewProtoWriter(cwr),
	}
}

// Read rewrite the method to use CountingReader
func (cc *CounterStream) Read(p []byte) (n int, err error) {
	return cc.rd.Read(p)
}

// Write rewrite the method to use CountingWriter
func (cc *CounterStream) Write(p []byte) (n int, err error) {
	return cc.wr.Write(p)
}

func (cc *CounterStream) ReadMessage(message proto.Message) error {
	return cc.ProtoMessageReader.ReadMessage(message)
}

func (cc *CounterStream) WriteMessage(m proto.Message) error {
	return cc.ProtoMessageWriter.WriteMessage(m)
}
