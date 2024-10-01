package connection

import (
	"net"
	"time"
)

func SetTcpOptions(conn net.Conn) (err error) {
	tcpConn := conn.(*net.TCPConn)
	
	if err = tcpConn.SetLinger(0); err != nil {
		return err
	}
	if err = tcpConn.SetNoDelay(false); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlivePeriod(60 * time.Second); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlive(true); err != nil {
		return err
	}
	
	return
}
