package dialer

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
)

var errUnexpectedInterfaceType = errors.New("unexpected interface type")

func InitDialer() {
	initTcpDialer()
	initQuicDialer()
}

func DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer, ok := proxy.FromEnvironment().(proxy.ContextDialer)
	if !ok {
		return nil, errUnexpectedInterfaceType
	}

	return dialer.DialContext(ctx, network, addr)
}

type CommonDialer struct {
	Local   *url.URL
	TlsConf *tls.Config
}
