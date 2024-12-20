package connection

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/pkg/common"

	"github.com/quic-go/quic-go"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

var (
	QuicConfig = &quic.Config{
		MaxIdleTimeout:  5 * time.Second,
		KeepAlivePeriod: 15 * time.Second,
		// MaxStreamReceiveWindow:     2 * 1024 * 1024,  // 2 MiB
		// MaxConnectionReceiveWindow: 10 * 1024 * 1024, // 10 MiB
	}
)

var _ connection.StreamConnection = (*QuicConnection)(nil)

type QuicConnection struct {
	connId   string
	udpConn  net.PacketConn
	quicConn quic.Connection
}

func NewQuicConnection(connId string, udpConn net.PacketConn, quicConn quic.Connection) *QuicConnection {
	return &QuicConnection{
		connId:   connId,
		udpConn:  udpConn,
		quicConn: quicConn,
	}
}

func (c *QuicConnection) GetConnectionId() string {
	return c.connId
}

func (c *QuicConnection) Close() error {
	if err := c.quicConn.CloseWithError(0, ""); err != nil {
		return err
	}

	if c.udpConn != nil {
		return c.udpConn.Close()
	}

	return nil
}

func (c *QuicConnection) RemoteAddr() net.Addr {
	return c.quicConn.RemoteAddr()
}

func (c *QuicConnection) LocalAddr() net.Addr {
	return c.quicConn.LocalAddr()
}

func (c *QuicConnection) SetDeadline(_ time.Time) error {
	return nil
}

func (c *QuicConnection) SetReadDeadline(_ time.Time) error {
	return nil
}

func (c *QuicConnection) SetWriteDeadline(_ time.Time) error {
	return nil
}

func (c *QuicConnection) AcceptStream() (connection.Stream, error) {
	stream, err := c.quicConn.AcceptStream(context.TODO())
	if err != nil {
		var idleTimeoutError *quic.IdleTimeoutError
		if errors.As(err, &idleTimeoutError) {
			return nil, errorutils.ErrConnectionRemoteDead
		}
		return nil, errorutils.WrapRemoteXError(err, "quicAcceptStream")
	}

	s := NewQuicStream(c.connId, c.LocalAddr(), c.RemoteAddr(), stream)
	s = connection.PackStream(s)

	return s, nil
}

func (c *QuicConnection) OpenStream() (connection.Stream, error) {
	stream, err := c.quicConn.OpenStream()
	if err != nil {
		return nil, err
	}
	s := NewQuicStream(c.connId, c.LocalAddr(), c.RemoteAddr(), stream)
	s = connection.PackStream(s)

	return s, nil
}

func (c *QuicConnection) ConnectionState() tls.ConnectionState {
	return c.quicConn.ConnectionState().TLS
}

var _ connection.Stream = (*QuicStream)(nil)

type QuicStream struct {
	quic.Stream
	protoutils.ProtoMessageReadWriter

	remoteAddr   net.Addr
	localAddr    net.Addr
	nodeId       common.NodeID
	connectionId string
}

func NewQuicStream(connId string, localAddr, remoteAddr net.Addr, stream quic.Stream) connection.Stream {
	return &QuicStream{
		Stream:                 stream,
		ProtoMessageReadWriter: NewStreamReadWriter(stream),

		remoteAddr:   remoteAddr,
		localAddr:    localAddr,
		connectionId: connId,
	}
}

func (q *QuicStream) GetConnectionId() string {
	return q.connectionId
}

func (q *QuicStream) GetNodeId() common.NodeID {
	return q.nodeId
}

func (q *QuicStream) SetNodeId(nodeId common.NodeID) {
	q.nodeId = nodeId
}

func (q *QuicStream) RemoteAddr() net.Addr {
	return q.remoteAddr
}

func (q *QuicStream) LocalAddr() net.Addr {
	return q.localAddr
}

func (q *QuicStream) Close() error {
	q.Stream.CancelRead(0)
	q.Stream.CancelWrite(0)
	return nil
}

func QuicNetwork(uri *url.URL) string {
	switch uri.Scheme {
	case "quic4":
		return "udp4"
	case "quic6":
		return "udp6"
	default:
		return "udp"
	}
}
