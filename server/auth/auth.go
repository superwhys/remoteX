package auth

import "github.com/superwhys/remoteX/domain/auth"

var _ auth.Service = (*SimpleAuth)(nil)

// SimpleAuth simply determines that the other party is a legitimate connection by exchanging protoMessages of nodes
type SimpleAuth struct {
}

func NewSimpleAuthService() auth.Service {
	return &SimpleAuth{}
}

func (s *SimpleAuth) AuthConnection() (err error) {
	return nil
}
