package auth

// Service is the auth service interface
// TODO: I haven't determined the authentication scheme yet. Let's write a empty method first
type Service interface {
	AuthConnection() (err error)
}

var _ Service = (*SimpleAuth)(nil)

// SimpleAuth simply determines that the other party is a legitimate connection by exchanging protoMessages of nodes
type SimpleAuth struct {
}

func NewSimpleAuthService() Service {
	return &SimpleAuth{}
}

func (s *SimpleAuth) AuthConnection() (err error) {
	return nil
}
