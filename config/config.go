package config

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
	
	m.LocalNode.Address = address
}

func (m *Config) Validate() error {
	return nil
}
