package api

import (
	"fmt"

	"github.com/superwhys/remoteX/domain/node"
)

type NodeDto struct {
	NodeId       string `json:"node_id"`
	ConnectionId string `json:"connection_id,omitempty"`
	Address      string `json:"address"`
	Status       string `json:"status"`
	IsLocal      bool   `json:"is_local"`
	Os           string `json:"os"`
	Arch         string `json:"arch"`
}

func NodeToDto(n *node.Node) *NodeDto {
	return &NodeDto{
		NodeId:       n.NodeId.String(),
		ConnectionId: n.ConnectionId,
		Address:      fmt.Sprintf("%s:%d", n.Address.GetIpAddress(), n.Address.GetPort()),
		Status:       n.Status.ToString(),
		IsLocal:      n.IsLocal,
		Os:           n.GetConfiguration().GetOs().String(),
		Arch:         n.GetConfiguration().GetArch().String(),
	}
}

func NodesToDto(nodes []*node.Node) []*NodeDto {
	var ns []*NodeDto
	for _, n := range nodes {
		ns = append(ns, NodeToDto(n))
	}
	return ns
}
