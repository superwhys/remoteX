// File:		tunnel.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package tunnel

import (
	"context"
	"fmt"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/command/tunnel"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protoutils"

	tunnelClient "github.com/superwhys/remoteX/internal/tunnel"
)

type ServiceImpl struct {
	tunnelClient *tunnelClient.TunnelManager
}

type tunnelArgs struct {
	localAddr  string
	remoteAddr string
}

func NewTunnelService() tunnel.Service {
	return &ServiceImpl{
		tunnelClient: tunnelClient.NewTunnelManager(),
	}
}

func (s *ServiceImpl) Name() string {
	return "tunnel"
}

func (s *ServiceImpl) SupportedCommands() []command.CommandType {
	return []command.CommandType{
		command.Forward,
		command.Forwardreceive,
		command.Listtunnel,
		command.Closetunnel,
	}
}

func (s *ServiceImpl) Invoke(ctx context.Context, cmd *command.Command, cmdCtx *command.CommandContext) (proto.Message, error) {
	switch cmd.Type {
	case command.Forward:
		return s.Forward(ctx, cmd.GetArgs(), cmdCtx.RawRemote)
	case command.Forwardreceive:
		return s.ReceiveForward(ctx, cmd.GetArgs(), cmdCtx.Remote)
	case command.Listtunnel:
		return s.ListTunnel(ctx, cmd.GetArgs())
	case command.Closetunnel:
		return s.CloseTunnel(ctx, cmd.GetArgs())
	default:
		return nil, errorutils.ErrCommandTypeNotSupport(cmd.GetType().String())
	}
}

func (s *ServiceImpl) ParseArgs(args command.Args) (*tunnelArgs, error) {
	var ea errorutils.ErrorArr
	a := &tunnelArgs{}
	for key, val := range args {
		switch key {
		case "local_addr":
			a.localAddr = val.GetStrValue()
		case "remote_addr":
			a.remoteAddr = val.GetStrValue()
		default:
			ea = append(ea, fmt.Errorf("error argument: %v", key))
		}
	}

	if len(ea) == 0 {
		return a, nil
	}

	return nil, errorutils.WrapRemoteXError(ea, "parseArgs")
}

func (s *ServiceImpl) validateTunnelArgs(tunnelArgs *tunnelArgs) error {
	if tunnelArgs.localAddr == "" || tunnelArgs.remoteAddr == "" {
		return errorutils.WrapRemoteXError(nil, "local_addr and remote_addr are required")
	}
	return nil
}

func (s *ServiceImpl) getStream(rw protoutils.ProtoMessageReadWriter) (connection.Stream, error) {
	stream, ok := rw.(connection.Stream)
	if !ok {
		return nil, errorutils.WrapRemoteXError(nil, "failed to convert rw to stream")
	}
	return stream, nil
}

func (s *ServiceImpl) Forward(ctx context.Context, args command.Args, conn connection.StreamConnection) (proto.Message, error) {
	tunnelArgs, err := s.ParseArgs(args)
	if err != nil {
		return nil, err
	}

	if err := s.validateTunnelArgs(tunnelArgs); err != nil {
		return nil, err
	}

	ctx = plog.With(ctx, "forward")
	tunnel, err := s.tunnelClient.CreateTunnel(ctx, tunnelArgs.localAddr, tunnelArgs.remoteAddr, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create forward tunnel: %v", err)
	}

	return &command.Tunnel{
		TunnelKey:  tunnel.TunnelKey,
		LocalAddr:  tunnel.LocalAddr,
		RemoteAddr: tunnel.RemoteAddr,
		Direction:  command.DirectionForward,
	}, nil
}

func (s *ServiceImpl) ReceiveForward(ctx context.Context, args command.Args, stream connection.Stream) (proto.Message, error) {
	tunnelKey, exists := args["tunnel_key"]
	if !exists {
		return nil, errorutils.ErrCommandMissingArguments(command.Reversereceive.String(), "tunnel_key")
	}
	localAddr, exists := args["addr"]
	if !exists {
		return nil, errorutils.ErrCommandMissingArguments(command.Reversereceive.String(), "addr")
	}

	plog.Debugc(ctx, "received forward request: %s, %s", tunnelKey.GetStrValue(), localAddr.GetStrValue())

	ctx = plog.With(ctx, "receiveForward")
	if err := s.tunnelClient.ReceiveTunnel(ctx, tunnelKey.GetStrValue(), localAddr.GetStrValue(), stream); err != nil {
		return nil, errorutils.WrapRemoteXError(err, "receiveForwardTunnel")
	}

	return nil, nil
}

func (s *ServiceImpl) ListTunnel(ctx context.Context, args command.Args) (proto.Message, error) {
	ts := s.tunnelClient.ListTunnels()

	resp := &command.ListTunnelResp{}
	for _, t := range ts {
		resp.Tunnels = append(resp.Tunnels, &command.Tunnel{
			TunnelKey:  t.TunnelKey,
			LocalAddr:  t.LocalAddr,
			RemoteAddr: t.RemoteAddr,
			Direction:  command.DirectionForward,
		})
	}

	return resp, nil
}

func (s *ServiceImpl) CloseTunnel(ctx context.Context, args command.Args) (proto.Message, error) {
	tunnelKey, exists := args["tunnel_key"]
	if !exists {
		return nil, errorutils.ErrCommandMissingArguments(command.Closetunnel.String(), "tunnel_key")
	}

	s.tunnelClient.CloseTunnel(tunnelKey.GetStrValue())

	return nil, nil
}
