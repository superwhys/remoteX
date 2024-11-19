package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/common"
)

type syncRequest struct {
	Target string `json:"target" binding:"required"`
	Src    string `json:"src" binding:"required"`
	Dest   string `json:"dest" binding:"required"`
	DryRun bool   `json:"dry_run"`
	Whole  bool   `json:"whole"`
}

func (p *syncRequest) toCommand(t command.CommandType) *command.Command {
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

func (a *RemoteXAPI) pullEntry(c *gin.Context, req *syncRequest) {
	callback := func(ctx context.Context, stream connection.Stream) error {
		localCmd := req.toCommand(command.Pull)
		_, err := a.srv.HandleCommand(ctx, localCmd, stream)
		if err != nil {
			return errors.Wrapf(err, "localCmd pull(%s -> %s)", req.Src, req.Dest)
		}
		return nil
	}

	remoteCmd := req.toCommand(command.Push)
	resp, err := a.srv.HandleRemoteCommand(c, common.NodeID(req.Target), remoteCmd, callback)
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("remote handle push command error: %v", err))
		return
	}

	pgin.ReturnSuccess(c, resp)
}

func (a *RemoteXAPI) pushEntry(c *gin.Context, req *syncRequest) {
	var (
		resp *command.Ret
		err  error
	)
	callback := func(ctx context.Context, stream connection.Stream) error {
		localCmd := req.toCommand(command.Push)
		resp, err = a.srv.HandleCommand(ctx, localCmd, stream)
		if err != nil {
			return errors.Wrapf(err, "localCmd push(%s -> %s)", req.Src, req.Dest)
		}
		return nil
	}

	remoteCmd := req.toCommand(command.Pull)
	_, err = a.srv.HandleRemoteCommand(c, common.NodeID(req.Target), remoteCmd, callback)
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("remote handle pull command error: %v", err))
		return
	}

	pgin.ReturnSuccess(c, resp)
}
