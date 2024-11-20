// File:		tunnel.go
// Created by:	Hoven
// Created on:	2024-11-20
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/pkg/common"
)

type TunnelRequest struct {
	TargetNode string `json:"target_node" binding:"required"`
	LocalAddr  string `json:"local_addr" binding:"required"`
	RemoteAddr string `json:"remote_addr" binding:"required"`
}

func (tr *TunnelRequest) toCommand(t command.CommandType) *command.Command {
	return &command.Command{
		Type: t,
		Args: map[string]command.Command_Arg{
			"local_addr":  command.StrArg(tr.LocalAddr),
			"remote_addr": command.StrArg(tr.RemoteAddr),
		},
	}
}

type TunnelResponse struct {
	Direction command.TunnelDirection `json:"direction"`
	TunnelKey string                  `json:"tunnel_key" `
}

type CloseTunnelRequest struct {
	TunnelKey string `json:"tunnel_key" binding:"required"`
}

func (a *RemoteXAPI) tunnelForward(c *gin.Context, req *TunnelRequest) {
	resp, err := a.srv.HandleCommandInBackground(c, common.NodeID(req.TargetNode), req.toCommand(command.Forward))
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("tunnel forward error: %v", err))
		return
	}

	pgin.ReturnSuccess(c, resp)
}

func (a *RemoteXAPI) listTunnel(c *gin.Context) {
	cmd := &command.Command{Type: command.Listtunnel}
	resp, err := a.srv.HandleLocalCommand(c, cmd)
	if err != nil {
		pgin.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("tunnel reverse error: %v", err))
		return
	}

	pgin.ReturnSuccess(c, resp)
}

func (a *RemoteXAPI) closeTunnel(ctx *gin.Context, req *CloseTunnelRequest) {
	cmd := &command.Command{
		Type: command.Closetunnel,
		Args: map[string]command.Command_Arg{
			"tunnel_key": command.StrArg(req.TunnelKey),
		},
	}
	_, err := a.srv.HandleLocalCommand(ctx, cmd)
	if err != nil {
		pgin.ReturnError(ctx, http.StatusInternalServerError, fmt.Sprintf("close tunnel error: %v", err))
		return
	}

	pgin.ReturnSuccess(ctx, nil)
}
