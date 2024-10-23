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
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/command/sync"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

var _ sync.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
}

func NewSyncService() *ServiceImpl {
	return &ServiceImpl{}
}

type argsSliceError []error

func (a argsSliceError) Error() string {
	if len(a) == 0 {
		return "no error"
	}

	s := ""
	for _, e := range a {
		s += e.Error() + "\n"
	}
	return s
}

func (s *ServiceImpl) ParseArgs(args command.Args) (opts *pb.SyncOpts, err error) {
	opts = opts.SetDefault()
	var se argsSliceError

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
			se = append(se, fmt.Errorf("error argument: %v", key))
		}
	}

	if len(se) == 0 {
		return opts, nil
	}

	return nil, se
}

func (s *ServiceImpl) Push(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error) {
	opts, err := s.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	if opts.Path == "" {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Push), args)
	}
	resp, err := filesync.SendFiles(ctx, rw, opts.Path, opts)
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
	if err := filesync.ReceiveFile(ctx, rw, opts.Dest, opts); err != nil {
		return nil, errorutils.ErrCommand(int32(command.Pull), args, errorutils.WithError(err))
	}

	return nil, nil
}
