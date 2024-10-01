package dialer

import (
	"context"
	"crypto/tls"
	"net/url"
	"time"
	
	"github.com/superwhys/remoteX/domain/connection"
	connPkg "github.com/superwhys/remoteX/internal/connection"
	"github.com/superwhys/remoteX/pkg/protocol"
	"github.com/superwhys/remoteX/pkg/tlsutils"
)

func initTcpDialer() {
	factory := &TcpDialerFactory{}
	for _, scheme := range []string{"tcp", "tcp4", "tcp6"} {
		connection.RegisterDialerFactory(scheme, factory)
	}
}

type tcpDialer struct {
	CommonDialer
}

func (t *tcpDialer) Dial(ctx context.Context, target *url.URL) (connection.TlsConn, error) {
	toCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	
	conn, err := DialContext(toCtx, target.Scheme, target.Host)
	if err != nil {
		return nil, err
	}
	
	if err := connPkg.SetTcpOptions(conn); err != nil {
		return nil, err
	}
	
	tc := tls.Client(conn, t.TlsConf)
	if err = tlsutils.TlsTimedHandshake(tc); err != nil {
		tc.Close()
		return nil, err
	}
	
	connId := connection.GenerateConnectionID(t.Local.Host, target.Host)
	
	sc, err := connPkg.NewTcpConnectionClient(connId, tc)
	if err != nil {
		tc.Close()
		return nil, err
	}
	
	c := &connection.Connection{
		ConnectionId:  connId,
		LocalAddress:  tc.LocalAddr().String(),
		RemoteAddress: tc.RemoteAddr().String(),
		Protocol:      protocol.ConnectionProtocolTcp,
		ConnectType:   protocol.ConnectionTypeClient,
		Status:        protocol.ConnectionStatusBeforeAuth,
		StartTime:     time.Now().Unix(),
		LastHeartbeat: time.Now().Unix(),
	}
	
	return connection.NewInternalConnection(sc, c, false), nil
}

type TcpDialerFactory struct{}

func (f *TcpDialerFactory) New(local *url.URL, tlsConf *tls.Config) connection.GenericDialer {
	return &tcpDialer{
		CommonDialer: CommonDialer{
			TlsConf: tlsConf,
			Local:   local,
		},
	}
}
