package node

import (
	"fmt"
	"net/url"
	"time"
)

func GetOsName(o string) NodeOS {
	switch o {
	case "linux":
		return NodeOsLinux
	case "windows":
		return NodeOsWin
	case "darwin":
		return NodeOsDarwin
	default:
		return NodeOsUnknown
	}
}

func GetArch(a string) NodeArch {
	switch a {
	case "amd64":
		return NodeArchX86
	case "arm64":
		return NodeArchArm
	default:
		return NodeArchUnknown
	}
}

func (m *Node) URL() *url.URL {
	return m.Address.URL()
}

func (m *Node) Host() string {
	return fmt.Sprintf("%s://%s:%d", m.Address.Schema, m.Address.IpAddress, m.Address.Port)
}

func (m *Node) CheckHeartbeatAlive() bool {
	now := time.Now().Unix()

	return now-m.GetLastHeartbeat() > int64(time.Second*20)
}

func (m *Node) CheckNodeOnline() bool {
	return m.Status == NodeStatusOnline && m.CheckHeartbeatAlive()
}

func (m *NodeTransConfiguration) SetDefault() {
	if m.MaxSendKbps == 0 {
		m.MaxSendKbps = 1024
	}
	if m.MaxRecvKbps == 0 {
		m.MaxRecvKbps = 1024
	}
}
