package listener

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"time"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/protocol"
	"github.com/superwhys/remoteX/pkg/tlsutils"
	
	connPkg "github.com/superwhys/remoteX/internal/connection"
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
				return ctx.Err()
			}
			return err
		default:
		}
		
		if err != nil {
			opErr := new(net.OpError)
			if errors.As(err, &opErr) && !opErr.Timeout() {
				plog.Warnf("Listen Tcp Accepting connection error:", err)
				
				acceptFailures++
				if acceptFailures > maxAcceptFailures {
					return err
				}
				
				time.Sleep(time.Duration(acceptFailures) * time.Second)
				
			}
			continue
		}
		
		acceptFailures = 0
		plog.Debugf("Listen TCP: connect from %v", conn.RemoteAddr())
		
		if err := connPkg.SetTcpOptions(conn); err != nil {
			plog.Errorf("SetTcpOptions error: %v", err)
			conn.Close()
			return err
		}
		tc := tls.Server(conn, tlsConf)
		if err := tlsutils.SetTimedHandshake(tc); err != nil {
			plog.Errorf("Listen TCP TLS handshake error: %v", err)
			tc.Close()
			continue
		}
		
		connId := connection.GenerateConnectionID(tlsutils.LocalHost(tc), tlsutils.RemoteHost(tc))
		sc, err := connPkg.NewTcpConnectionServer(connId, tc)
		if err != nil {
			plog.Errorf("new tcp connection server error: %v", err)
			tc.Close()
			continue
		}
		
		c := &connection.Connection{
			ConnectionId:  connId,
			LocalAddress:  tc.LocalAddr().String(),
			RemoteAddress: tc.RemoteAddr().String(),
			Protocol:      protocol.ConnectionProtocolTcp,
			ConnectType:   protocol.ConnectionTypeServer,
			Status:        protocol.ConnectionStatusBeforeAuth,
			StartTime:     time.Now().Unix(),
			LastHeartbeat: time.Now().Unix(),
		}
		ic := connection.NewInternalConnection(sc, c, t.local, true)
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
