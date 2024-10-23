// File:		sync.go
// Created by:	Hoven
// Created on:	2024-10-18
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package sync

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/command/sync"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

var _ sync.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
}

func NewSyncService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s *ServiceImpl) Push(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error) {
	path := args["path"]
	if path == "" {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Push), args)
	}
	resp, err := filesync.SendFiles(ctx, rw, path, nil)
	if err != nil {
		return nil, errorutils.ErrCommand(int32(command.Push), args, errorutils.WithError(err))
	}

	return resp, nil
}

func (s *ServiceImpl) Pull(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error) {
	dest := args["dest"]
	if dest == "" {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Pull), args)
	}
	if err := filesync.ReceiveFile(ctx, rw, dest, nil); err != nil {
		return nil, errorutils.ErrCommand(int32(command.Pull), args, errorutils.WithError(err))
	}

	return nil, nil
}
