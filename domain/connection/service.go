package connection

import (
	"context"
	"crypto/tls"
	"net/url"
	
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/pkg/common"
	errors2 "github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protocol"
)

// Service is the interface of Connection domain
// The application layer will run these methods in the following order
// 1. CreateListener in suture.Service
// 2. HandleListenerConnection
type Service interface {
	// CreateListener will create corresponding Listeners based on the configuration of the Node
	// and listen for connections from other Nodes in a separate coroutine
	// and then transmit the TlsConn to the outer layer through a channel
	CreateListener(ctx context.Context, conns chan<- TlsConn) error
	// EstablishConnection will create a connection with target Node in the role of a client
	EstablishConnection(ctx context.Context, target *url.URL) (TlsConn, error)
	CheckConnection(conn TlsConn) error
	CloseConnection(connId string) error
}

var _ Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	localNodeId string
	local       *url.URL
	tlsConf     *tls.Config
	connections map[string]TlsConn
}

func NewConnectionService(local *url.URL, tlsConf *tls.Config) Service {
	return &ServiceImpl{
		local:       local,
		tlsConf:     tlsConf,
		connections: make(map[string]TlsConn),
	}
}

func (s *ServiceImpl) CreateListener(ctx context.Context, conns chan<- TlsConn) error {
	creator, err := GetListenerFactory(s.local)
	if err != nil {
		return err
	}
	lis := creator.New(s.local, s.tlsConf)
	
	return lis.Listen(ctx, conns)
}

func (s *ServiceImpl) EstablishConnection(ctx context.Context, target *url.URL) (TlsConn, error) {
	dialFactory, err := GetDialerFactory(target)
	if err != nil {
		return nil, err
	}
	
	tlsConn, err := dialFactory.New(s.local, s.tlsConf).Dial(ctx, target)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to establish connection to %v", target)
	}
	
	connectionId := GenerateConnectionID(s.local.Host, target.Host)
	s.connections[connectionId] = tlsConn
	
	return tlsConn, nil
}

func (s *ServiceImpl) CheckConnection(conn TlsConn) error {
	cs := conn.ConnectionState()
	certs := cs.PeerCertificates
	if cl := len(certs); cl != 1 {
		return errors2.ErrConnection(conn, errors2.WithMsg("peer certificate invalidate"))
	}
	remoteCert := certs[0]
	remoteID := common.NewNodeID(remoteCert.Raw)
	if remoteID.String() == s.localNodeId {
		return errors2.ErrConnectToMyself(remoteID, conn)
	}
	
	return nil
}

func (s *ServiceImpl) CloseConnection(connId string) error {
	conn, ok := s.connections[connId]
	if !ok {
		return errors.New("连接未找到")
	}
	conn.SetStatus(protocol.ConnectionStatusDisconnected)
	delete(s.connections, connId)
	return conn.Close()
}

func (s *ServiceImpl) SendMessage(connId string, message proto.Message) error {
	conn, ok := s.connections[connId]
	if !ok {
		return errors.New("连接未找到")
	}
	
	conn.UpdateLastHeartbeat()
	return conn.WriteMessage(message)
}
