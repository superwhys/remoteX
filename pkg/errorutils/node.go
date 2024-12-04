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

func ErrNodeNotFound(nodeId common.NodeID) *NodeError {
	return &NodeError{
		BaseError: PackBaseError(WithMsg("node not found")),
		nodeId:    nodeId,
	}
}
