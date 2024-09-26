package dialer

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"time"
	
	"github.com/superwhys/remoteX/domain/connection"
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
	
	if err := t.setTcpOptions(conn); err != nil {
		return nil, err
	}
	
	tc := tls.Client(conn, t.TlsConf)
	if err = tlsutils.TlsTimedHandshake(tc); err != nil {
		tc.Close()
		return nil, err
	}
	
	connId := connection.GenerateConnectionID(t.Local.Host, target.Host)
	c := &connection.Connection{
		ConnectionId:  connId,
		Protocol:      protocol.ConnectionProtocolTcp,
		ConnectType:   protocol.ConnectionTypeClient,
		Status:        protocol.ConnectionStatusBeforeAuth,
		StartTime:     time.Now().Unix(),
		LastHeartbeat: time.Now().Unix(),
	}
	
	return connection.NewInternalConn(tc, c), nil
}

func (t *tcpDialer) setTcpOptions(conn net.Conn) (err error) {
	tcpConn := conn.(*net.TCPConn)
	
	if err = tcpConn.SetLinger(0); err != nil {
		return err
	}
	if err = tcpConn.SetNoDelay(false); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlivePeriod(60 * time.Second); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlive(true); err != nil {
		return err
	}
	
	return
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
