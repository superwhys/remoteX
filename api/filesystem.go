package api

import (
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/pkg/common"
)

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

func (a *RemoteXAPI) listDir(c *gin.Context, req *listDirReq) {
	cmd := req.toCommand(command.Listdir)
	ret, err := a.srv.HandleCommand(c, cmd, nil)
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

func (a *RemoteXAPI) listRemoteDir(c *gin.Context, req *listRemoteDir) {
	cmd := req.toCommand(command.Listdir)
	resp, err := a.srv.HandleRemoteCommand(c, common.NodeID(req.NodeId), cmd, nil)
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("handle remote command error: %v", err))
		return
	}
	
	pgin.ReturnSuccess(c, resp)
}
