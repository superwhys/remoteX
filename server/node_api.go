package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/pkg/common"
)

func (s *RemoteXServer) getAllNodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		nodes, err := s.nodeService.GetNodes()
		if err != nil {
			pgin.ReturnError(c, http.StatusBadRequest, err.Error())
			return
		}

		pgin.ReturnSuccess(c, nodes)
	}
}

type getNode struct {
	NodeId string `uri:"nodeId"`
}

func (s *RemoteXServer) getNode(c *gin.Context, req *getNode) {
	plog.Infof("get node request: %v", req)
	node, err := s.nodeService.GetNode(common.NodeID(req.NodeId))
	if err != nil {
		pgin.ReturnError(c, http.StatusNotFound, err.Error())
		return
	}

	pgin.ReturnSuccess(c, node)
}
