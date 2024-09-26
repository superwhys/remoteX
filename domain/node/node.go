package node

import (
	"fmt"
	"net/url"
)

func (m *Node) URL() *url.URL {
	return m.Address.URL()
}

func (m *Node) Host() string {
	return fmt.Sprintf("%s://%s:%d", m.Address.Schema, m.Address.IpAddress, m.Address.Port)
}
