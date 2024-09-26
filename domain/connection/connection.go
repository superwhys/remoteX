package connection

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	fmt "fmt"
	"io"
	"net"
	"sync"
	"time"
	
	"github.com/superwhys/remoteX/pkg/protocol"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type TlsConn interface {
	io.ReadWriteCloser
	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter
	ConnectionState() tls.ConnectionState
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	SetDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
	GetConnectionId() string
	GetStatus() protocol.ConnectionStatus
	GetProtocol() protocol.ConnectionProtocol
	GetConnectType() protocol.ConnectionType
	GetStartTime() int64
	SetStatus(status protocol.ConnectionStatus)
	UpdateLastHeartbeat()
	String() string
}

var _ TlsConn = (*InternalConnection)(nil)

type InternalConnection struct {
	sync.Mutex
	*tls.Conn
	*Connection
	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter
	
	TlsConf      *tls.Config
	connectionID string
}

func NewInternalConn(tc *tls.Conn, conn *Connection) *InternalConnection {
	return &InternalConnection{
		Conn:               tc,
		ProtoMessageReader: protoutils.NewProtoReader(tc),
		ProtoMessageWriter: protoutils.NewProtoWriter(tc),
		Connection:         conn,
	}
}

func (c *InternalConnection) SetStatus(status protocol.ConnectionStatus) {
	c.Status = status
}

func (c *InternalConnection) UpdateLastHeartbeat() {
	c.LastHeartbeat = time.Now().Unix()
}

func (c *InternalConnection) String() string {
	return fmt.Sprintf("%s-%s/%s", c.LocalAddr(), c.RemoteAddr(), c.GetConnectionId())
}

func GenerateConnectionID(source, target string) string {
	rawID := fmt.Sprintf("%s-%s/%d", source, target, time.Now().Unix())
	
	hash := sha256.New()
	hash.Write([]byte(rawID))
	
	return hex.EncodeToString(hash.Sum(nil))
}
