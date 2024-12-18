package connection

import (
	"context"
	"crypto/tls"
	"net/url"

	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protocol"
)

var _ connection.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	localNodeId string
	local       *url.URL
	tlsConf     *tls.Config
	connections map[string]connection.TlsConn
}

func NewConnectionService(local *url.URL, tlsConf *tls.Config) connection.Service {
	return &ServiceImpl{
		local:       local,
		tlsConf:     tlsConf,
		connections: make(map[string]connection.TlsConn),
	}
}

func (s *ServiceImpl) CreateListener(ctx context.Context, connCh chan<- connection.TlsConn) error {
	creator, err := connection.GetListenerFactory(s.local)
	if err != nil {
		return errorutils.ErrGetListenerCreator(err)
	}
	lis := creator.New(s.local, s.tlsConf)

	return lis.Listen(ctx, connCh)
}

func (s *ServiceImpl) EstablishConnection(ctx context.Context, target *url.URL) (connection.TlsConn, error) {
	creator, err := connection.GetDialerFactory(target)
	if err != nil {
		return nil, errorutils.ErrGetDialerCreator(err)
	}

	streamConn, err := creator.New(s.local, s.tlsConf).Dial(ctx, target)
	if err != nil {
		return nil, errorutils.ErrEstablishConnection(err)
	}

	return streamConn, nil
}

func (s *ServiceImpl) CheckConnection(conn connection.TlsConn) error {
	cs := conn.ConnectionState()
	certs := cs.PeerCertificates
	if cl := len(certs); cl != 1 {
		return errorutils.ErrConnectionCert
	}

	remoteCert := certs[0]
	remoteID := common.NewNodeID(remoteCert.Raw)
	if remoteID.String() == s.localNodeId {
		return errorutils.ErrConnectToMyself
	}

	return nil
}

func (s *ServiceImpl) RegisterConnection(conn connection.TlsConn) {
	s.connections[conn.GetConnectionId()] = conn
}

func (s *ServiceImpl) GetConnection(connId string) (connection.TlsConn, error) {
	conn, ok := s.connections[connId]
	if !ok {
		return nil, errorutils.ErrConnectNotFound
	}
	return conn, nil
}

func (s *ServiceImpl) CloseConnection(connId string) error {
	conn, ok := s.connections[connId]
	if !ok {
		return errorutils.ErrConnectNotFound
	}
	conn.SetStatus(protocol.ConnectionStatusDisconnected)
	delete(s.connections, connId)
	return conn.Close()
}
