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
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/command/sync"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/pkg/errorutils"
)

var _ sync.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
}

func NewSyncService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s *ServiceImpl) ParseArgs(args command.Args) (opts *pb.SyncOpts, err error) {
	opts = opts.SetDefault()
	var ea errorutils.ErrorArr

	for key, val := range args {
		switch key {
		case "path":
			opts.Path = val.GetStrValue()
		case "dest":
			opts.Dest = val.GetStrValue()
		case "dry_run":
			opts.DryRun = val.GetBoolValue()
		case "whole":
			opts.Whole = val.GetBoolValue()
		default:
			ea = append(ea, fmt.Errorf("error argument: %v", key))
		}
	}

	if len(ea) == 0 {
		return opts, nil
	}

	return nil, ea
}

func (s *ServiceImpl) Invoke(ctx context.Context, cmd *command.Command, opt *command.RemoteOpt) (proto.Message, error) {
	switch cmd.Type {
	case command.Push:
		return s.Push(ctx, cmd.Args, opt)
	case command.Pull:
		return s.Pull(ctx, cmd.Args, opt)
	default:
		return nil, errorutils.ErrNotSupportCommandType(int32(cmd.Type))
	}
}

func (s *ServiceImpl) tellRemoteTo(remoteType command.CommandType, currentArgs command.Args, stream connection.Stream) error {
	remoteCmd := &command.Command{
		Type: remoteType,
		Args: currentArgs,
	}

	if err := stream.WriteMessage(remoteCmd); err != nil {
		return errors.Wrap(err, "sendCommand")
	}

	return nil
}

func (s *ServiceImpl) Push(ctx context.Context, args command.Args, opt *command.RemoteOpt) (proto.Message, error) {
	opts, err := s.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	if opts.Path == "" {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Push), args)
	}

	stream, err := opt.Conn.OpenStream()
	if err != nil {
		return nil, errorutils.ErrCommand(int32(command.Push), args, errorutils.WithError(err))
	}

	err = s.tellRemoteTo(command.Pull, args, stream)
	if err != nil {
		return nil, errorutils.ErrCommand(int32(command.Push), args, errorutils.WithError(err))
	}

	resp, err := filesync.SendFiles(ctx, stream, opts.Path, opts)
	if err != nil {
		return nil, errorutils.ErrCommand(int32(command.Push), args, errorutils.WithError(err))
	}

	return resp, nil
}

func (s *ServiceImpl) Pull(ctx context.Context, args command.Args, opt *command.RemoteOpt) (proto.Message, error) {
	opts, err := s.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	if opts.Dest == "" {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Pull), args)
	}

	stream, err := opt.Conn.OpenStream()
	if err != nil {
		return nil, errorutils.ErrCommand(int32(command.Push), args, errorutils.WithError(err))
	}

	err = s.tellRemoteTo(command.Push, args, stream)
	if err != nil {
		return nil, errorutils.ErrCommand(int32(command.Pull), args, errorutils.WithError(err))
	}

	if err := filesync.ReceiveFile(ctx, stream, opts.Dest, opts); err != nil {
		return nil, errorutils.ErrCommand(int32(command.Pull), args, errorutils.WithError(err))
	}

	return nil, nil
}
