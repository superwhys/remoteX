package auth

import (
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
)

type Service interface {
	AuthConnection(remote *node.Node, conn connection.TlsConn) (err error)
}

var _ Service = (*SimpleAuth)(nil)

// SimpleAuth simply determines that the other party is a legitimate connection by exchanging protoMessages of nodes
type SimpleAuth struct {
	local *node.Node
}

func NewSimpleAuthService(local *node.Node) Service {
	return &SimpleAuth{
		local: local,
	}
}

func (s *SimpleAuth) AuthConnection(remote *node.Node, conn connection.TlsConn) (err error) {
	return nil
}
