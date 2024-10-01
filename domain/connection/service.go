package connection

import (
	"context"
	"net/url"
)

// Service is the interface of Connection domain
// The application layer will run these methods in the following order
// 1. CreateListener in suture.Service
// 2. HandleListenerConnection
type Service interface {
	// CreateListener will create corresponding Listeners based on the configuration of the Node
	// and listen for connections from other Nodes in a separate coroutine
	// and then transmit the TlsConn to the outer layer through a channel
	CreateListener(ctx context.Context, connCh chan<- TlsConn) error
	// EstablishConnection will create a connection with target Node in the role of a client
	EstablishConnection(ctx context.Context, target *url.URL) (TlsConn, error)
	CheckConnection(conn TlsConn) error
	RegisterConnection(conn TlsConn)
	GetConnection(connId string) (TlsConn, error)
	CloseConnection(connId string) error
}
