package listener

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/quic-go/quic-go"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/protocol"

	connPkg "github.com/superwhys/remoteX/internal/connection"
)

func initQuicListener() {
	factory := &QuicListenerFactory{}
	for _, scheme := range []string{"quic", "quic4", "quic6"} {
		connection.RegisterListenerFactory(scheme, factory)
	}
}

var _ connection.GenericListener = (*QuicListener)(nil)

type QuicListener struct {
	local   *url.URL
	tlsConf *tls.Config
}

func (l *QuicListener) Listen(ctx context.Context, connections chan<- connection.TlsConn) error {
	network := connPkg.QuicNetwork(l.local)

	udpAddr, err := net.ResolveUDPAddr(network, l.local.Host)
	if err != nil {
		plog.Errorc(ctx, "Resolve quic addr error: %v", err)
		return err
	}

	udpConn, err := net.ListenUDP(network, udpAddr)
	if err != nil {
		plog.Errorc(ctx, "Listen quic addr error: %v", err)
		return err
	}
	defer udpConn.Close()

	quicTransport := &quic.Transport{
		Conn: udpConn,
	}
	defer quicTransport.Close()

	listener, err := quicTransport.Listen(l.tlsConf, connPkg.QuicConfig)
	if err != nil {
		return err
	}
	defer listener.Close()

	plog.Infoc(ctx, "Quic listener (%v) starting", udpAddr)
	defer plog.Infoc(ctx, "Quic listener (%v) shutting down", udpAddr)

	acceptFailures := 0
	const maxAcceptFailures = 10

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		conn, err := listener.Accept(ctx)
		if errors.Is(err, context.Canceled) {
			return nil
		} else if err != nil {
			plog.Warnf("Listen Quic Accepting connection error:", err)

			acceptFailures++
			if acceptFailures > maxAcceptFailures {
				return err
			}

			time.Sleep(time.Duration(acceptFailures) * time.Second)
			continue
		}

		plog.Debugf("Listen Quic: connect from %v", conn.RemoteAddr())

		connId := connection.GenerateConnectionID(conn.LocalAddr().String(), conn.RemoteAddr().String())
		sc := connPkg.NewQuicConnection(connId, nil, conn)

		c := &connection.Connection{
			ConnectionId:  connId,
			LocalAddress:  conn.LocalAddr().String(),
			RemoteAddress: conn.RemoteAddr().String(),
			Protocol:      protocol.ConnectionProtocolTcp,
			ConnectType:   protocol.ConnectionTypeServer,
			Status:        protocol.ConnectionStatusBeforeAuth,
			StartTime:     time.Now().Unix(),
			LastHeartbeat: time.Now().Unix(),
		}
		ic := connection.NewInternalConnection(sc, c, l.local, true)
		connections <- ic
	}
}

type QuicListenerFactory struct{}

func (t *QuicListenerFactory) New(local *url.URL, tlsConf *tls.Config) connection.GenericListener {
	l := &QuicListener{
		local:   local,
		tlsConf: tlsConf,
	}

	return l
}
