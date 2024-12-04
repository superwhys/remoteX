package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/server"
)

type listDir struct {
	Path string `form:"path" binding:"required"`
}

func (l *listDir) toCommand(t command.CommandType) *command.Command {
	return &command.Command{
		Type: t,
		Args: map[string]command.Command_Arg{
			"path": command.StrArg(l.Path),
		},
	}
}

func (a listDir) Handle(c *gin.Context, srv *server.RemoteXServer) (resp *command.Ret, err error) {
	cmd := a.toCommand(command.Listdir)
	ret, err := srv.HandleLocalCommand(c, cmd)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type listRemoteDir struct {
	NodeId string `uri:"nodeId" binding:"required"`
	Path   string `form:"path" binding:"required"`
}

func (ld *listRemoteDir) toCommand(t command.CommandType) *command.Command {
	return &command.Command{
		Type: t,
		Args: map[string]command.Command_Arg{
			"path": command.StrArg(ld.Path),
		},
	}
}

func (ld listRemoteDir) Handle(c *gin.Context, srv *server.RemoteXServer) (resp *command.Ret, err error) {
	cmd := ld.toCommand(command.Listdir)
	resp, err = srv.HandleCommandWithRemote(c, common.NodeID(ld.NodeId), cmd)
	if err != nil {
		return nil, errors.Wrap(err, "list remote dir")
	}

	return resp, nil
}
