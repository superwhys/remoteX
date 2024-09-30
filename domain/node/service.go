package node

import (
	"errors"
	"sync"
	"time"

	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/osutils"
)

type Service interface {
	RegisterNode(n *Node) error
	GetNode(nodeId common.NodeID) (*Node, error)
	GetLocal() *Node
	RefreshCurrentNode() (*Node, error)
	GetNodeStatus(nodeId common.NodeID) (NodeStatus, error)
	UpdateNodeStatus(nodeId common.NodeID, status NodeStatus) error
	UpdateHeartbeat(nodeId common.NodeID) error
}

type ServiceImpl struct {
	localNode *Node
	rl        *sync.RWMutex
	nodes     map[common.NodeID]*Node
}

func NewNodeService(local *Node) *ServiceImpl {
	return &ServiceImpl{
		localNode: local,
		rl:        &sync.RWMutex{},
		nodes:     make(map[common.NodeID]*Node),
	}
}

func (ds *ServiceImpl) RegisterNode(n *Node) error {
	if _, exists := ds.nodes[n.NodeId]; exists {
		return errors.New("该设备已注册")
	}

	addr := n.Address
	if addr.IpAddress == "" || addr.Port == 0 {
		return errors.New("无效的设备信息: 缺少 IP 地址或端口")
	}
	ds.nodes[n.NodeId] = n
	return nil
}

func (ds *ServiceImpl) GetNode(nodeId common.NodeID) (*Node, error) {
	node, exists := ds.nodes[nodeId]
	if !exists {
		return nil, errorutils.ErrNodeNotFound(nodeId)
	}

	return node, nil
}

func (ds *ServiceImpl) GetLocal() *Node {
	ds.rl.RLock()
	defer ds.rl.RUnlock()

	return ds.localNode
}

func (ds *ServiceImpl) RefreshCurrentNode() (*Node, error) {
	ds.rl.RLock()
	currentNode := ds.localNode
	ds.rl.RUnlock()

	os, arch := osutils.GetOsArch()
	currentNode.Configuration.Os = GetOsName(os)
	currentNode.Configuration.Arch = GetArch(arch)
	currentNode.LastHeartbeat = time.Now().Unix()

	ds.rl.Lock()
	defer ds.rl.Unlock()

	ds.localNode = currentNode
	return currentNode, nil
}

func (ds *ServiceImpl) UpdateNodeStatus(nodeId common.NodeID, status NodeStatus) error {
	n, ok := ds.nodes[nodeId]
	if !ok {
		return errors.New("设备未找到")
	}
	n.Status = status
	return nil
}

func (ds *ServiceImpl) GetNodeStatus(nodeId common.NodeID) (NodeStatus, error) {
	n, ok := ds.nodes[nodeId]
	if !ok {
		return NodeStatus(0), errors.New("设备未找到")
	}
	return n.Status, nil
}

func (ds *ServiceImpl) UpdateHeartbeat(nodeId common.NodeID) error {
	n, ok := ds.nodes[nodeId]
	if !ok {
		return errors.New("设备未找到")
	}
	n.LastHeartbeat = time.Now().Unix()
	return nil
}
