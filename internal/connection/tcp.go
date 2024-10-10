package connection

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/protoutils"
	"github.com/xtaci/smux"
)

var _ connection.StreamConnection = (*TcpConnection)(nil)

type TcpConnection struct {
	session *smux.Session
	rawConn *tls.Conn
	connId  string
}

func (c *TcpConnection) ConnectionState() tls.ConnectionState {
	return c.rawConn.ConnectionState()
}

func NewTcpConnectionServer(connId string, conn *tls.Conn) (*TcpConnection, error) {
	session, err := smux.Server(conn, nil)
	if err != nil {
		return nil, err
	}

	return &TcpConnection{
		session: session,
		rawConn: conn,
		connId:  connId,
	}, nil
}

func NewTcpConnectionClient(connId string, conn *tls.Conn) (*TcpConnection, error) {
	session, err := smux.Client(conn, nil)
	if err != nil {
		return nil, err
	}

	return &TcpConnection{
		session: session,
		rawConn: conn,
		connId:  connId,
	}, nil
}

func (c *TcpConnection) RemoteAddr() net.Addr {
	return c.rawConn.RemoteAddr()
}

func (c *TcpConnection) LocalAddr() net.Addr {
	return c.rawConn.LocalAddr()
}

func (c *TcpConnection) SetDeadline(t time.Time) error {
	return c.session.SetDeadline(t)
}

func (c *TcpConnection) SetReadDeadline(t time.Time) error {
	return c.SetDeadline(t)
}

func (c *TcpConnection) SetWriteDeadline(t time.Time) error {
	return c.SetDeadline(t)
}

func (c *TcpConnection) AcceptStream() (connection.Stream, error) {
	s, err := c.session.AcceptStream()
	if err != nil {
		return nil, err
	}

	return NewTcpStream(c.connId, s), nil
}

func (c *TcpConnection) OpenStream() (connection.Stream, error) {
	s, err := c.session.OpenStream()
	if err != nil {
		return nil, err
	}

	return NewTcpStream(c.connId, s), nil
}

func (c *TcpConnection) Close() error {
	c.session.Close()
	return c.rawConn.Close()
}

var _ connection.Stream = (*TcpStream)(nil)

type TcpStream struct {
	*smux.Stream
	*StreamReadWriter

	pr           *protoutils.ProtoReader
	pw           *protoutils.ProtoWriter
	nodeId       common.NodeID
	connectionId string
}

func NewTcpStream(connId string, stream *smux.Stream) *TcpStream {
	return &TcpStream{
		Stream: stream,
		StreamReadWriter: &StreamReadWriter{
			pr:           protoutils.NewProtoReader(stream),
			pw:           protoutils.NewProtoWriter(stream),
			connectionId: connId,
		},
	}
}
