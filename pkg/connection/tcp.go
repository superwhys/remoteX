package connection

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/protoutils"
	"github.com/xtaci/smux"
)

var _ connection.StreamConnection = (*TcpConnection)(nil)

type TcpConnection struct {
	session *smux.Session
	rawConn *tls.Conn
}

func (c *TcpConnection) ConnectionState() tls.ConnectionState {
	return c.rawConn.ConnectionState()
}

func NewTcpConnectionServer(conn *tls.Conn) *TcpConnection {
	session, err := smux.Server(conn, nil)
	if err != nil {
		panic(err)
	}

	return &TcpConnection{
		session: session,
		rawConn: conn,
	}
}

func NewTcpConnectionClient(conn *tls.Conn) *TcpConnection {
	session, err := smux.Client(conn, nil)
	if err != nil {
		panic(err)
	}

	return &TcpConnection{
		session: session,
		rawConn: conn,
	}
}

func (c *TcpConnection) RemoteAddr() net.Addr {
	return c.rawConn.RemoteAddr()
}

func (c *TcpConnection) LocalAddr() net.Addr {
	return c.rawConn.LocalAddr()
}

func (c *TcpConnection) SetDeadline(t time.Time) error {
	fmt.Println("tcpConnSetDeadline")
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

	return NewTcpStream(s), nil
}

func (c *TcpConnection) OpenStream() (connection.Stream, error) {
	s, err := c.session.OpenStream()
	if err != nil {
		return nil, err
	}

	return NewTcpStream(s), nil
}

func (c *TcpConnection) Close() error {
	c.session.Close()
	return c.rawConn.Close()
}

var _ connection.Stream = (*TcpStream)(nil)

type TcpStream struct {
	*smux.Stream

	pr     *protoutils.ProtoReader
	pw     *protoutils.ProtoWriter
	nodeId common.NodeID
}

func NewTcpStream(stream *smux.Stream) *TcpStream {
	return &TcpStream{
		Stream: stream,
		pr:     protoutils.NewProtoReader(stream),
		pw:     protoutils.NewProtoWriter(stream),
	}
}

func (t *TcpStream) ReadMessage(message proto.Message) error {
	return t.pr.ReadMessage(message)
}

func (t *TcpStream) WriteMessage(m proto.Message) error {
	return t.pw.WriteMessage(m)
}

func (t *TcpStream) GetNodeId() common.NodeID {
	return t.nodeId
}

func (t *TcpStream) SetNodeId(nodeId common.NodeID) {
	t.nodeId = nodeId
}
