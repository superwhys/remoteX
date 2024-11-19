package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/superwhys/remoteX/pkg/common"
)

func (a *RemoteXAPI) getAllNodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		nodes, err := a.srv.NodeService.GetNodes()
		if err != nil {
			pgin.ReturnError(c, http.StatusBadRequest, err.Error())
			return
		}

		pgin.ReturnSuccess(c, NodesToDto(nodes))
	}
}

type getNode struct {
	NodeId string `uri:"nodeId" binding:"required"`
}

func (a *RemoteXAPI) getNode(c *gin.Context, req *getNode) {
	node, err := a.srv.NodeService.GetNode(common.NodeID(req.NodeId))
	if err != nil {
		pgin.ReturnError(c, http.StatusNotFound, err.Error())
		return
	}

	pgin.ReturnSuccess(c, NodeToDto(node))
}
