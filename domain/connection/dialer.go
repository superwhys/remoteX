package connection

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
)

var (
	dialerFactories = make(map[string]DialerFactory)
)

type GenericDialer interface {
	Dial(ctx context.Context, target *url.URL) (TlsConn, error)
}

type DialerFactory interface {
	New(local *url.URL, tlsConf *tls.Config) GenericDialer
}

func RegisterDialerFactory(schema string, factory DialerFactory) {
	dialerFactories[schema] = factory
}

func GetDialerFactory(uri *url.URL) (DialerFactory, error) {
	dialerFactory, ok := dialerFactories[uri.Scheme]
	if !ok {
		return nil, fmt.Errorf("unknown address scheme %q", uri.Scheme)
	}

	return dialerFactory, nil
}
