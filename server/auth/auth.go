package auth

import (
	"github.com/superwhys/remoteX/domain/auth"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
)

type AuthAppService struct {
	service auth.Service
}

func NewAuthAppService(local *node.Node) *AuthAppService {
	return &AuthAppService{
		service: auth.NewSimpleAuthService(local),
	}
}

func (a *AuthAppService) AuthRemoteConn(remote *node.Node, conn connection.TlsConn) (err error) {
	return a.service.AuthConnection(remote, conn)
}
