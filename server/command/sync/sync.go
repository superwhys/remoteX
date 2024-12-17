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
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

var _ sync.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	filesyncer filesync.FileSyncer
}

func NewSyncService() *ServiceImpl {
	return &ServiceImpl{
		filesyncer: filesync.NewFileSyncer(),
	}
}

func (s *ServiceImpl) Name() string {
	return "sync"
}

func (s *ServiceImpl) SupportedCommands() []command.CommandType {
	return []command.CommandType{command.Push, command.Pull}
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

func (s *ServiceImpl) Invoke(ctx context.Context, cmd *command.Command, cmdCtx *command.CommandContext) (proto.Message, error) {
	if cmdCtx.IsRemote {
		err := s.tellRemoteTo(cmd, cmdCtx.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "tellRemoteTo")
		}
	}

	switch cmd.Type {
	case command.Push:
		return s.Push(ctx, cmd.Args, cmdCtx.Remote)
	case command.Pull:
		return s.Pull(ctx, cmd.Args, cmdCtx.Remote)
	default:
		return nil, errorutils.ErrNotSupportCommandType(int32(cmd.Type))
	}
}

func (s *ServiceImpl) tellRemoteTo(cmd *command.Command, rw protoutils.ProtoMessageReadWriter) error {
	var remoteCmdType command.CommandType
	switch cmd.GetType() {
	case command.Push:
		remoteCmdType = command.Pull
	case command.Pull:
		remoteCmdType = command.Push
	default:
		return errorutils.ErrNotSupportCommandType(int32(cmd.GetType()))
	}

	remoteCmd := &command.Command{
		Type: remoteCmdType,
		Args: cmd.GetArgs(),
	}

	if err := rw.WriteMessage(remoteCmd); err != nil {
		return errors.Wrap(err, "sendCommand")
	}

	return nil
}

func (s *ServiceImpl) Push(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error) {
	opts, err := s.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	if opts.Path == "" {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Push), args)
	}

	resp, err := s.filesyncer.SendFiles(ctx, rw, opts.Path, opts)
	if err != nil {
		return nil, errorutils.ErrCommand(int32(command.Push), args, errorutils.WithError(err))
	}

	return resp, nil
}

func (s *ServiceImpl) Pull(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error) {
	opts, err := s.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	if opts.Dest == "" {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Pull), args)
	}

	resp, err := s.filesyncer.ReceiveFiles(ctx, rw, opts.Dest, opts)
	if err != nil {
		return nil, errorutils.ErrCommand(int32(command.Pull), args, errorutils.WithError(err))
	}

	return resp, nil
}
