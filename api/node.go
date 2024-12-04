package api

import (
	"github.com/gin-gonic/gin"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/server"
)

type getAllNodes struct{}

func (g getAllNodes) Handle(c *gin.Context, srv *server.RemoteXServer) (resp []*NodeDto, err error) {
	nodes, err := srv.NodeService.GetNodes()
	if err != nil {
		return nil, err
	}

	return NodesToDto(nodes), nil
}

type getNode struct {
	NodeId string `uri:"nodeId" binding:"required"`
}

func (a getNode) Handle(c *gin.Context, srv *server.RemoteXServer) (resp *NodeDto, err error) {
	node, err := srv.NodeService.GetNode(common.NodeID(a.NodeId))
	if err != nil {
		return nil, err
	}

	return NodeToDto(node), nil
}
