package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/pkg/common"
)

func (s *RemoteXServer) getAllNodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		nodes, err := s.nodeService.GetNodes()
		if err != nil {
			pgin.ReturnError(c, http.StatusBadRequest, err.Error())
			return
		}

		pgin.ReturnSuccess(c, NodesToDtos(nodes))
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

	pgin.ReturnSuccess(c, NodeToDto(node))
}

type listDirReq struct {
	Path string `form:"path"`
}

func (s *RemoteXServer) listDir(c *gin.Context, req *listDirReq) {
	cmd := &command.Command{Type: command.Listdir, Args: map[string]string{"path": req.Path}}
	ret, err := s.commandService.DoCommand(c, cmd)
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pgin.ReturnSuccess(c, ret)
}

type listRemoteDir struct {
	NodeId string `uri:"nodeId"`
	Path   string `form:"path"`
}

func (s *RemoteXServer) listRemoteDir(c *gin.Context, req *listRemoteDir) {
	if req.NodeId == "" {
		pgin.ReturnError(c, http.StatusBadRequest, "remote nodeId is required")
		return
	}
	cmd := &command.Command{Type: command.Listdir, Args: map[string]string{"path": req.Path}}

	resp, err := s.handleRemoteCommand(c, common.NodeID(req.NodeId), cmd)
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("handle remote command error: %v", err))
		return
	}

	pgin.ReturnSuccess(c, resp)
}
