package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/server"
)

type sync struct {
	Target string `json:"target" binding:"required"`
	Src    string `json:"src" binding:"required"`
	Dest   string `json:"dest" binding:"required"`
	DryRun bool   `json:"dry_run"`
	Whole  bool   `json:"whole"`
}

func (p *sync) toCommand(t command.CommandType) *command.Command {
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

type syncPull struct {
	sync
}

func (p syncPull) Handle(c *gin.Context, srv *server.RemoteXServer) (resp *command.Ret, err error) {
	cmd := p.toCommand(command.Pull)
	resp, err = srv.HandleCommandWithRemote(c, common.NodeID(p.Target), cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "localCmd pull(%s -> %s)", p.Src, p.Dest)
	}

	return resp, nil
}

type syncPush struct {
	sync
}

func (p syncPush) Handle(c *gin.Context, srv *server.RemoteXServer) (resp *command.Ret, err error) {
	cmd := p.toCommand(command.Push)
	resp, err = srv.HandleCommandWithRemote(c, common.NodeID(p.Target), cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "localCmd push(%s -> %s)", p.Src, p.Dest)
	}

	return resp, nil
}
