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

func (l *listDirReq) toCommand(t command.CommandType) *command.Command {
	return &command.Command{
		Type: t,
		Args: map[string]command.Command_Arg{
			"path": command.StrArg(l.Path),
		},
	}
}

func (s *RemoteXServer) listDir(c *gin.Context, req *listDirReq) {
	cmd := req.toCommand(command.Listdir)
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

func (l *listRemoteDir) toCommand(t command.CommandType) *command.Command {
	return &command.Command{
		Type: t,
		Args: map[string]command.Command_Arg{
			"path": command.StrArg(l.Path),
		},
	}
}

func (s *RemoteXServer) listRemoteDir(c *gin.Context, req *listRemoteDir) {
	cmd := req.toCommand(command.Listdir)
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
	DryRun bool   `json:"dry_run"`
	Whole  bool   `json:"whole"`
}

func (p *pullEntry) toCommand(t command.CommandType) *command.Command {
	return &command.Command{
		Type: t,
		Args: map[string]command.Command_Arg{
			"path":    command.StrArg(p.Src),
			"dest":    command.StrArg(p.Dest),
			"dry_run": command.BoolArg(p.DryRun),
			"whole":   command.BoolArg(p.Whole),
		},
	}
}

func (s *RemoteXServer) pullEntry(c *gin.Context, req *pullEntry) {
	callback := func(ctx context.Context, stream connection.Stream) error {
		localCmd := req.toCommand(command.Pull)
		_, err := s.commandService.DoCommand(ctx, localCmd, stream)
		return errors.Wrapf(err, "localCmd pull(%s -> %s)", req.Src, req.Dest)
	}

	remoteCmd := req.toCommand(command.Push)
	resp, err := s.handleRemoteCommand(c, common.NodeID(req.Target), remoteCmd, callback)
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("handle remote command error: %v", err))
		return
	}

	pgin.ReturnSuccess(c, resp)
}
