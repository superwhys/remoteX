package dialer

import (
	"context"
	"crypto/tls"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/protocol"

	connPkg "github.com/superwhys/remoteX/internal/connection"
	tlsutils "github.com/superwhys/remoteX/internal/tls"
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
		return nil, errors.Wrap(err, "dial target")
	}

	if err := connPkg.SetTcpOptions(conn); err != nil {
		return nil, errors.Wrap(err, "set tcp options")
	}

	tc := tls.Client(conn, t.TlsConf)
	if err = tlsutils.TlsTimedHandshake(tc); err != nil {
		tc.Close()
		return nil, errors.Wrap(err, "tls handshake")
	}

	connId := connection.GenerateConnectionID(t.Local.Host, target.Host)

	sc, err := connPkg.NewTcpConnectionClient(connId, tc)
	if err != nil {
		tc.Close()
		return nil, errors.Wrap(err, "new tcp connection client")
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

	return connection.NewInternalConnection(sc, c, target, false), nil
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
