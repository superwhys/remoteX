package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/pgin"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
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
	NodeId string `uri:"nodeId" binding:"required"`
}

func (s *RemoteXServer) getNode(c *gin.Context, req *getNode) {
	node, err := s.nodeService.GetNode(common.NodeID(req.NodeId))
	if err != nil {
		pgin.ReturnError(c, http.StatusNotFound, err.Error())
		return
	}

	pgin.ReturnSuccess(c, NodeToDto(node))
}

type listDirReq struct {
	Path string `form:"path" binding:"required"`
}

func (s *RemoteXServer) listDir(c *gin.Context, req *listDirReq) {
	cmd := &command.Command{Type: command.Listdir, Args: map[string]string{"path": req.Path}}
	ret, err := s.commandService.DoCommand(c, cmd, nil)
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pgin.ReturnSuccess(c, ret)
}

type listRemoteDir struct {
	NodeId string `uri:"nodeId" binding:"required"`
	Path   string `form:"path" binding:"required"`
}

func (s *RemoteXServer) listRemoteDir(c *gin.Context, req *listRemoteDir) {
	cmd := &command.Command{Type: command.Listdir, Args: map[string]string{"path": req.Path}}

	resp, err := s.handleRemoteCommand(c, common.NodeID(req.NodeId), cmd, nil)
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("handle remote command error: %v", err))
		return
	}

	pgin.ReturnSuccess(c, resp)
}

type pullEntry struct {
	Target string `json:"target" binding:"required"`
	Src    string `json:"src" binding:"required"`
	Dest   string `json:"dest" binding:"required"`
}

func (s *RemoteXServer) pullEntry(c *gin.Context, req *pullEntry) {
	remoteCmd := &command.Command{Type: command.Push, Args: map[string]string{"path": req.Src}}

	resp, err := s.handleRemoteCommand(c, common.NodeID(req.Target), remoteCmd, func(ctx context.Context, stream connection.Stream) error {
		localCmd := &command.Command{Type: command.Pull, Args: map[string]string{"dest": req.Dest}}
		_, err := s.commandService.DoCommand(ctx, localCmd, stream)
		return errors.Wrapf(err, "localCmd pull(%s -> %s)", req.Src, req.Dest)
	})
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("handle remote command error: %v", err))
		return
	}

	pgin.ReturnSuccess(c, resp)
}
