// File:		proxy.go
// Created by:	Hoven
// Created on:	2024-12-21
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/server"
)

type proxy struct {
	Node string `header:"X-Node" binding:"required"`
	Addr string `header:"X-Remote-Addr" binding:"tcp_addr"`
}

func (p proxy) toCommand() *command.Command {
	return &command.Command{
		Type: command.Proxy,
		Args: map[string]command.Command_Arg{
			"is_server":  command.BoolArg(true),
			"proxy_addr": command.StrArg(p.Addr),
		},
	}
}

func (p proxy) Handle(c *gin.Context, srv *server.RemoteXServer) (resp *command.Ret, err error) {
	err = srv.HandleCommandWithReqExchange(context.TODO(), common.NodeID(p.Node), p.toCommand(), c.Request, c.Writer)
	if err != nil {
		return nil, errorutils.WrapRemoteXError(err, "handleCommandWithWriterExchange")
	}

	return nil, nil
}
