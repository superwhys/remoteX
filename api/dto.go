package api

import (
	"fmt"
	
	"github.com/superwhys/remoteX/domain/node"
)

type NodeDto struct {
	NodeId       string `json:"node_id"`
	ConnectionId string `json:"connection_id,omitempty"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Status       string `json:"status"`
	IsLocal      bool   `json:"is_local"`
}

func NodeToDto(n *node.Node) *NodeDto {
	return &NodeDto{
		NodeId:       n.NodeId.String(),
		ConnectionId: n.ConnectionId,
		Name:         n.Name,
		Address:      fmt.Sprintf("%s:%d", n.Address.GetIpAddress(), n.Address.GetPort()),
		Status:       n.Status.ToString(),
		IsLocal:      n.IsLocal,
	}
}

func NodesToDto(nodes []*node.Node) []*NodeDto {
	var ns []*NodeDto
	for _, n := range nodes {
		ns = append(ns, NodeToDto(n))
	}
	return ns
}
