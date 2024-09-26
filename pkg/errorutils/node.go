package errorutils

import (
	"fmt"
	
	"github.com/superwhys/remoteX/pkg/common"
)

type NodeError struct {
	*BaseError
	nodeId common.NodeID
}

func ErrNode(nodeId common.NodeID, opts ...ErrOption) *NodeError {
	ne := &NodeError{
		BaseError: &BaseError{},
		nodeId:    nodeId,
	}
	
	for _, opt := range opts {
		opt(ne.BaseError)
	}
	
	return ne
}

func (ne *NodeError) String() string {
	return fmt.Sprintf("NodeError{%v}. Error{%v}", ne.nodeId, ne.BaseError)
}

func (ne *NodeError) Error() string {
	return ne.String()
}

type NodeNotFound struct {
	*NodeError
}

func ErrNodeNotFound(nodeId common.NodeID) *NodeNotFound {
	return &NodeNotFound{
		NodeError: ErrNode(nodeId),
	}
}
