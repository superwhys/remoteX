package connection

import (
	"crypto/tls"
	"io"
	"math"
	"net"
	"net/url"
	"time"

	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/protocol"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type DialTask struct {
	Target *url.URL
	NodeId common.NodeID

	DialCnt   int
	Threshold int
	InitDelay time.Duration
	MaxDelay  time.Duration
	IsRedial  bool
}

func calculateThreshold(initialDelay, maxDelay time.Duration) int {
	return int(math.Log2(float64(maxDelay) / float64(initialDelay)))
}

func NewDialTask(target *url.URL, nodeId common.NodeID, initDelay, maxDelay time.Duration, isRedial bool) *DialTask {
	return &DialTask{
		Target:    target,
		NodeId:    nodeId,
		InitDelay: initDelay,
		MaxDelay:  maxDelay,
		Threshold: calculateThreshold(initDelay, maxDelay),
		IsRedial:  isRedial,
	}
}

type BaseConnection interface {
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}

type ConnInformation interface {
	GetConnectionId() string
}

type NodeInformation interface {
	GetNodeId() common.NodeID
	SetNodeId(nodeId common.NodeID)
}

type StreamConnection interface {
	io.Closer
	BaseConnection
	ConnInformation

	AcceptStream() (Stream, error)
	OpenStream() (Stream, error)
	ConnectionState() tls.ConnectionState
}

type Stream interface {
	io.ReadWriteCloser
	BaseConnection
	ConnInformation
	NodeInformation

	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter
}

type TlsConn interface {
	StreamConnection
	NodeInformation

	IsServer() bool
	GetDialURL() *url.URL
	SetStatus(status protocol.ConnectionStatus)
	UpdateLastHeartbeat()
}

var _ TlsConn = (*InternalConn)(nil)

type InternalConn struct {
	StreamConnection
	*Connection
	nodeId   common.NodeID
	isServer bool
	target   *url.URL
}

func NewInternalConnection(sc StreamConnection, conn *Connection, target *url.URL, isServer bool) *InternalConn {
	return &InternalConn{
		StreamConnection: sc,
		Connection:       conn,
		isServer:         isServer,
		target:           target,
	}
}

func (i *InternalConn) GetNodeId() common.NodeID {
	return i.nodeId
}

func (i *InternalConn) SetNodeId(nodeId common.NodeID) {
	i.nodeId = nodeId
}

func (i *InternalConn) GetConnectionId() string {
	return i.StreamConnection.GetConnectionId()
}

func (i *InternalConn) IsServer() bool {
	return i.isServer
}

func (i *InternalConn) GetDialURL() *url.URL {
	return i.target
}

func (i *InternalConn) AcceptStream() (Stream, error) {
	stream, err := i.StreamConnection.AcceptStream()
	if err != nil {
		return nil, err
	}
	stream.SetNodeId(i.nodeId)
	return stream, nil
}

func (i *InternalConn) OpenStream() (Stream, error) {
	stream, err := i.StreamConnection.OpenStream()
	if err != nil {
		return nil, err
	}
	stream.SetNodeId(i.nodeId)

	return stream, nil
}

func (c *Connection) SetStatus(status protocol.ConnectionStatus) {
	c.Status = status
}

func (c *Connection) UpdateLastHeartbeat() {
	c.LastHeartbeat = time.Now().Unix()
}
