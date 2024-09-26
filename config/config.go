package config

import (
	"github.com/superwhys/remoteX/domain/node"
)

func (m *Config) SetDefault() {
	address := m.LocalNode.Address
	if address.IpAddress == "" {
		address.IpAddress = "127.0.0.1"
	}

	if address.Port == 0 {
		address.Port = 28080
	}

	if address.Schema == "" {
		address.Schema = "tcp"
	}

	if m.TransConf == nil {
		m.TransConf = &node.NodeTransConfiguration{}
		m.TransConf.SetDefault()
	}

	m.LocalNode.Address = address
	m.LocalNode.Configuration = &node.NodeConfiguration{Transmission: m.TransConf}
}

func (m *Config) Validate() error {
	return nil
}
