package node

import (
	"github.com/superwhys/remoteX/pkg/common"
)

type Service interface {
	RegisterNode(n *Node) error
	GetNode(nodeId common.NodeID) (*Node, error)
	GetNodes() ([]*Node, error)
	GetLocal() *Node
	RefreshCurrentNode() (*Node, error)
	GetNodeStatus(nodeId common.NodeID) (NodeStatus, error)
	UpdateNodeStatus(nodeId common.NodeID, status NodeStatus) error
	UpdateHeartbeat(nodeId common.NodeID) error
}
