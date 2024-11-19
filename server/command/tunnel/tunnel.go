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

	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/command/tunnel"
	"github.com/superwhys/remoteX/pkg/protoutils"

	tunnelClient "github.com/superwhys/remoteX/internal/tunnel"
)

var _ tunnel.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	tunnelClient *tunnelClient.SshTunnel
}

func NewTunnelService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s *ServiceImpl) PasreArgs(args command.Args)

func (s *ServiceImpl) Forward(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error) {
	return nil, nil
}

func (s *ServiceImpl) Reverse(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error) {
	return nil, nil
}
