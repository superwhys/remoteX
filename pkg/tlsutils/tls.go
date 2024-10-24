package tlsutils

import (
	"crypto/tls"
	"fmt"
	"time"
)

var (
	tlsHandshakeTimeout = 10 * time.Second
)

func SetTimedHandshake(tc *tls.Conn) error {
	tc.SetDeadline(time.Now().Add(tlsHandshakeTimeout))
	defer tc.SetDeadline(time.Time{})
	
	return tc.Handshake()
}

func LocalHost(tc *tls.Conn) string {
	return fmt.Sprintf("%s://%s", tc.LocalAddr().Network(), tc.LocalAddr().String())
}

func RemoteHost(tc *tls.Conn) string {
	return fmt.Sprintf("%s://%s", tc.RemoteAddr().Network(), tc.RemoteAddr().String())
}
