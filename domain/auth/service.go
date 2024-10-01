package auth

// Service is the auth service interface
// TODO: I haven't determined the authentication scheme yet. Let's write a empty method first
type Service interface {
	AuthConnection() (err error)
}
