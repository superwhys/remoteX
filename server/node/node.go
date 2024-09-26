package node

import (
	"github.com/superwhys/remoteX/domain/node"
	"github.com/superwhys/remoteX/pkg/protocol"
)

type NodeAppService struct {
	node.Service
}

func NewNodeAppService() *NodeAppService {
	return &NodeAppService{
		Service: node.NewNodeService(),
	}
}

func (ns *NodeAppService) NewNode() *node.Node {
	return &node.Node{
		NodeId:        "",
		ConnectionId:  "",
		Name:          "",
		Address:       protocol.Address{},
		Status:        0,
		Configuration: nil,
		LastHeartbeat: 0,
	}
}
