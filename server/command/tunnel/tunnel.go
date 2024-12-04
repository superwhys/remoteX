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

var _ tunnel.Service = (*ServiceImpl)(nil)

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

func (s *ServiceImpl) Invoke(ctx context.Context, cmd *command.Command, opt *command.RemoteOpt) (proto.Message, error) {
	switch cmd.Type {
	case command.Forward:
		return s.Forward(ctx, cmd.GetArgs(), opt)
	case command.Forwardreceive:
		return s.ReceiveForward(ctx, cmd.GetArgs(), opt)
	case command.Listtunnel:
		return s.ListTunnel(ctx, cmd.GetArgs(), opt)
	case command.Closetunnel:
		return s.CloseTunnel(ctx, cmd.GetArgs(), opt)
	default:
		return nil, errorutils.ErrNotSupportCommandType(int32(cmd.Type))
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

	return nil, ea
}

func (s *ServiceImpl) validateTunnelArgs(tunnelArgs *tunnelArgs) error {
	if tunnelArgs.localAddr == "" || tunnelArgs.remoteAddr == "" {
		return fmt.Errorf("local_addr and remote_addr are required")
	}
	return nil
}

func (s *ServiceImpl) getStream(rw protoutils.ProtoMessageReadWriter) (connection.Stream, error) {
	stream, ok := rw.(connection.Stream)
	if !ok {
		return nil, fmt.Errorf("failed to convert rw to stream")
	}
	return stream, nil
}

func (s *ServiceImpl) Forward(ctx context.Context, args command.Args, opt *command.RemoteOpt) (proto.Message, error) {
	if opt.Conn == nil {
		return nil, fmt.Errorf("connection is required")
	}
	tunnelArgs, err := s.ParseArgs(args)
	if err != nil {
		return nil, err
	}

	if err := s.validateTunnelArgs(tunnelArgs); err != nil {
		return nil, err
	}

	ctx = plog.With(ctx, "forward")
	tunnel, err := s.tunnelClient.CreateTunnel(ctx, tunnelArgs.localAddr, tunnelArgs.remoteAddr, opt.Conn)
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

func (s *ServiceImpl) ReceiveForward(ctx context.Context, args command.Args, opt *command.RemoteOpt) (proto.Message, error) {
	tunnelKey, exists := args["tunnel_key"]
	if !exists {
		return nil, fmt.Errorf("tunnel_key argument is required")
	}
	localAddr, exists := args["addr"]
	if !exists {
		return nil, fmt.Errorf("addr argument is required")
	}

	plog.Debugc(ctx, "received forward request: %s, %s", tunnelKey.GetStrValue(), localAddr.GetStrValue())

	ctx = plog.With(ctx, "receiveForward")
	if err := s.tunnelClient.ReceiveTunnel(ctx, tunnelKey.GetStrValue(), localAddr.GetStrValue(), opt.Stream); err != nil {
		return nil, fmt.Errorf("failed to receive forward tunnel: %v", err)
	}

	return nil, nil
}

func (s *ServiceImpl) ListTunnel(ctx context.Context, args command.Args, _ *command.RemoteOpt) (proto.Message, error) {
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

func (s *ServiceImpl) CloseTunnel(ctx context.Context, args command.Args, _ *command.RemoteOpt) (proto.Message, error) {
	tunnelKey, exists := args["tunnel_key"]
	if !exists {
		return nil, fmt.Errorf("tunnel_key argument is required")
	}

	s.tunnelClient.CloseTunnel(tunnelKey.GetStrValue())

	return nil, nil
}
