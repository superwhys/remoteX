package dialer

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/protocol"

	connPkg "github.com/superwhys/remoteX/internal/connection"
)

func initQuicDialer() {
	factory := &QuicDialerFactory{}
	for _, scheme := range []string{"quic", "quic4", "quic6"} {
		connection.RegisterDialerFactory(scheme, factory)
	}
}

var (
	quicDialTimeout = time.Second * 10
)

type quicDialer struct {
	CommonDialer
}

func (q *quicDialer) Dial(ctx context.Context, target *url.URL) (connection.TlsConn, error) {
	network := connPkg.QuicNetwork(target)

	addr, err := net.ResolveUDPAddr(network, target.Host)
	if err != nil {
		return nil, errors.Wrap(err, "resolve udp addr")
	}

	udpPacketConn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return nil, errors.Wrap(err, "listen udp packet")
	}
	transport := &quic.Transport{Conn: udpPacketConn}

	ctx, cancel := context.WithTimeout(ctx, quicDialTimeout)
	defer cancel()

	quicConn, err := transport.Dial(ctx, addr, q.TlsConf, connPkg.QuicConfig)
	if err != nil {
		udpPacketConn.Close()
		return nil, errors.Wrap(err, "dial quic transport")
	}

	connId := connection.GenerateConnectionID(q.Local.Host, target.Host)
	sc := connPkg.NewQuicConnection(connId, udpPacketConn, quicConn)
	c := &connection.Connection{
		ConnectionId:  connId,
		LocalAddress:  quicConn.LocalAddr().String(),
		RemoteAddress: quicConn.RemoteAddr().String(),
		Protocol:      protocol.ConnectionProtocolTcp,
		ConnectType:   protocol.ConnectionTypeClient,
		Status:        protocol.ConnectionStatusBeforeAuth,
		StartTime:     time.Now().Unix(),
		LastHeartbeat: time.Now().Unix(),
	}

	return connection.NewInternalConnection(sc, c, target, false), nil
}

type QuicDialerFactory struct{}

func (f *QuicDialerFactory) New(local *url.URL, tlsConf *tls.Config) connection.GenericDialer {
	return &quicDialer{
		CommonDialer: CommonDialer{
			TlsConf: tlsConf,
			Local:   local,
		},
	}
}
