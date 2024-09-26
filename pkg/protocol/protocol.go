package protocol

import (
	"fmt"
	"net/url"
)

func (m *Address) URL() *url.URL {
	return &url.URL{
		Scheme: m.Schema,
		Host:   fmt.Sprintf("%s:%d", m.IpAddress, m.Port),
	}
}
