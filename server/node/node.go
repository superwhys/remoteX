package node

import (
	"sync"
	"time"

	"github.com/superwhys/remoteX/domain/node"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/osutils"
)

type ServiceImpl struct {
	localNode *node.Node
	rl        *sync.RWMutex
	nodes     map[common.NodeID]*node.Node
}

func NewNodeService(local *node.Node) *ServiceImpl {
	s := &ServiceImpl{
		localNode: new(node.Node),
		rl:        &sync.RWMutex{},
		nodes:     make(map[common.NodeID]*node.Node),
	}
	copyNode := *local
	s.localNode = &copyNode

	s.RegisterNode(local)

	return s
}

func (ds *ServiceImpl) RegisterNode(n *node.Node) error {
	oldNode, _ := ds.nodes[n.NodeId]
	if oldNode != nil && oldNode.CheckNodeOnline() {
		return errorutils.ErrNode(n.NodeId, errorutils.WithMsg("Has a same online node"))
	}

	addr := n.Address
	if addr.IpAddress == "" || addr.Port == 0 {
		return errorutils.ErrNode(n.NodeId, errorutils.WithMsg("missing addr or port"))
	}
	ds.nodes[n.NodeId] = n
	return nil
}

func (ds *ServiceImpl) GetNode(nodeId common.NodeID) (*node.Node, error) {
	node, exists := ds.nodes[nodeId]
	if !exists {
		return nil, errorutils.ErrNodeNotFound(nodeId)
	}

	return node, nil
}

func (ds *ServiceImpl) GetNodes() ([]*node.Node, error) {
	ds.rl.RLock()
	defer ds.rl.RUnlock()

	nodes := make([]*node.Node, 0, len(ds.nodes))
	for _, n := range ds.nodes {
		nodes = append(nodes, n)
	}

	return nodes, nil
}

func (ds *ServiceImpl) GetLocal() *node.Node {
	ds.rl.RLock()
	defer ds.rl.RUnlock()

	return ds.localNode
}

func (ds *ServiceImpl) RefreshCurrentNode() (*node.Node, error) {
	ds.rl.RLock()
	currentNode := ds.localNode
	ds.rl.RUnlock()

	os, arch := osutils.GetOsArch()
	currentNode.Configuration.Os = node.GetOsName(os)
	currentNode.Configuration.Arch = node.GetArch(arch)
	currentNode.LastHeartbeat = time.Now().Unix()

	ds.rl.Lock()
	defer ds.rl.Unlock()

	ds.localNode = currentNode
	return currentNode, nil
}

func (ds *ServiceImpl) UpdateNode(n *node.Node) error {
	n, ok := ds.nodes[n.NodeId]
	if !ok {
		return errorutils.ErrNodeNotFound(n.NodeId)
	}

	ds.nodes[n.NodeId] = n
	return nil
}

func (ds *ServiceImpl) UpdateNodeStatus(nodeId common.NodeID, status node.NodeStatus) error {
	n, ok := ds.nodes[nodeId]
	if !ok {
		return errorutils.ErrNodeNotFound(nodeId)
	}
	n.Status = status
	if status == node.NodeStatusOffline {
		n.ConnectionId = ""
	}
	return nil
}

func (ds *ServiceImpl) GetNodeStatus(nodeId common.NodeID) (node.NodeStatus, error) {
	n, ok := ds.nodes[nodeId]
	if !ok {
		return node.NodeStatus(0), errorutils.ErrNodeNotFound(nodeId)
	}
	return n.Status, nil
}

func (ds *ServiceImpl) UpdateHeartbeat(nodeId common.NodeID) error {
	n, ok := ds.nodes[nodeId]
	if !ok {
		return errorutils.ErrNodeNotFound(nodeId)
	}
	n.LastHeartbeat = time.Now().Unix()
	return nil
}
