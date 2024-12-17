// File:		tunnel.go
// Created by:	Hoven
// Created on:	2024-11-20
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/server"
)

type tunnelForward struct {
	TargetNode string `json:"target_node" binding:"required"`
	LocalAddr  string `json:"local_addr" binding:"required"`
	RemoteAddr string `json:"remote_addr" binding:"required"`
}

func (tr *tunnelForward) toCommand(t command.CommandType) *command.Command {
	return &command.Command{
		Type: t,
		Args: map[string]command.Command_Arg{
			"local_addr":  command.StrArg(tr.LocalAddr),
			"remote_addr": command.StrArg(tr.RemoteAddr),
		},
	}
}

func (t tunnelForward) Handle(c *gin.Context, srv *server.RemoteXServer) (resp *command.Ret, err error) {
	resp, err = srv.HandleCommandWithRawRemote(context.TODO(), common.NodeID(t.TargetNode), t.toCommand(command.Forward))
	if err != nil {
		return nil, errors.Wrap(err, "tunnel forward")
	}

	return resp, nil
}

type listTunnel struct{}

func (lt listTunnel) Handle(c *gin.Context, srv *server.RemoteXServer) (resp *command.Ret, err error) {
	cmd := &command.Command{Type: command.Listtunnel}
	resp, err = srv.HandleLocalCommand(c, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "list tunnel")
	}

	return resp, nil
}

type closeTunnel struct {
	TunnelKey string `json:"tunnel_key" binding:"required"`
}

func (ct closeTunnel) Handle(c *gin.Context, srv *server.RemoteXServer) (any, err error) {
	cmd := &command.Command{
		Type: command.Closetunnel,
		Args: map[string]command.Command_Arg{
			"tunnel_key": command.StrArg(ct.TunnelKey),
		},
	}
	_, err = srv.HandleLocalCommand(c, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "closeTunnel")
	}

	return nil, nil
}
