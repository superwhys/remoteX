package listener

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"time"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/protocol"
	"github.com/superwhys/remoteX/pkg/tlsutils"
)

func initTcpListener() {
	factory := &TcpListenerFactory{}
	for _, scheme := range []string{"tcp", "tcp4", "tcp6"} {
		connection.RegisterListenerFactory(scheme, factory)
	}
}

type TcpListener struct {
	local   *url.URL
	tlsConf *tls.Config
}

func (t *TcpListener) Listen(ctx context.Context, connections chan<- connection.TlsConn) error {
	local := t.local
	tlsConf := t.tlsConf
	
	tcpAddr, err := net.ResolveTCPAddr(local.Scheme, local.Host)
	if err != nil {
		plog.Errorc(ctx, "Resolve tcp addr error: %v", err)
		return err
	}
	
	lisConf := net.ListenConfig{}
	
	listener, err := lisConf.Listen(ctx, local.Scheme, tcpAddr.String())
	if err != nil {
		plog.Errorc(ctx, "Listen tcp addr error: %v", err)
		return err
	}
	defer listener.Close()
	
	plog.Infoc(ctx, "TCP listener (%v) starting", tcpAddr)
	defer plog.Infoc(ctx, "TCP listener (%v) shutting down", tcpAddr)
	
	acceptFailures := 0
	const maxAcceptFailures = 10
	
	tcpListener := listener.(*net.TCPListener)
	for {
		_ = tcpListener.SetDeadline(time.Now().Add(time.Second))
		conn, err := tcpListener.Accept()
		select {
		case <-ctx.Done():
			if err == nil {
				conn.Close()
			}
			return nil
		default:
		}
		
		if err != nil {
			if err, ok := err.(*net.OpError); !ok || !err.Timeout() {
				plog.Warnc(ctx, "Listen Tcp Accepting connection error:", err)
				
				acceptFailures++
				if acceptFailures > maxAcceptFailures {
					return err
				}
				
				time.Sleep(time.Duration(acceptFailures) * time.Second)
			}
			continue
		}
		
		acceptFailures = 0
		plog.Debugc(ctx, "Listen TCP: connect from %v", conn.RemoteAddr())
		
		tc := tls.Server(conn, tlsConf)
		if err := tlsutils.TlsTimedHandshake(tc); err != nil {
			plog.Errorc(ctx, "Listen TCP TLS handshake error: %v", err)
			tc.Close()
			continue
		}
		
		c := &connection.Connection{
			ConnectionId:  connection.GenerateConnectionID(tlsutils.LocalHost(tc), tlsutils.RemoteHost(tc)),
			Protocol:      protocol.ConnectionProtocolTcp,
			ConnectType:   protocol.ConnectionTypeServer,
			Status:        protocol.ConnectionStatusBeforeAuth,
			StartTime:     time.Now().Unix(),
			LastHeartbeat: time.Now().Unix(),
		}
		ic := connection.NewInternalConn(tc, c)
		connections <- ic
	}
}

type TcpListenerFactory struct{}

func (t *TcpListenerFactory) New(local *url.URL, tlsConf *tls.Config) connection.GenericListener {
	l := &TcpListener{
		local:   local,
		tlsConf: tlsConf,
	}
	
	return l
}
