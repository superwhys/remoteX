package connection

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
)

var (
	listenerFactories = make(map[string]ListenerFactory)
)

type GenericListener interface {
	Listen(ctx context.Context, connCh chan<- TlsConn) error
}

type ListenerFactory interface {
	New(local *url.URL, tlsConf *tls.Config) GenericListener
}

func RegisterListenerFactory(schema string, factory ListenerFactory) {
	listenerFactories[schema] = factory
}

func GetListenerFactory(uri *url.URL) (ListenerFactory, error) {
	dialerFactory, ok := listenerFactories[uri.Scheme]
	if !ok {
		return nil, fmt.Errorf("unknown address scheme %q", uri.Scheme)
	}

	return dialerFactory, nil
}
