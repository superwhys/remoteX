// File:		service.go
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
	"github.com/superwhys/remoteX/domain/connection"
)

type Service interface {
	command.CommandProvider
	Forward(ctx context.Context, args command.Args, conn connection.StreamConnection) (proto.Message, error)
	ReceiveForward(ctx context.Context, args command.Args, stream connection.Stream) (proto.Message, error)
	ListTunnel(ctx context.Context, args command.Args) (proto.Message, error)
	CloseTunnel(ctx context.Context, args command.Args) (proto.Message, error)
}
