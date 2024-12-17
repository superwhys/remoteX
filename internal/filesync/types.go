// File:		types.go
// Created by:	Hoven
// Created on:	2024-12-16
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package filesync

import (
	"context"

	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/internal/filesync/receiver"
	"github.com/superwhys/remoteX/internal/filesync/sender"
	"github.com/superwhys/remoteX/internal/filesystem"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type FileSyncer interface {
	SendFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, path string, opts *pb.SyncOpts) (*pb.SyncResp, error)
	ReceiveFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, path string, opts *pb.SyncOpts) (*pb.SyncResp, error)
}

type Transfer interface {
	Transfer(ctx context.Context, file *pb.FileBase, opts *pb.SyncOpts) (err error)
	TransferSize() int64
}

func NewSendTransfer(rw protoutils.ProtoMessageReadWriter) Transfer {
	return &sender.SendTransfer{
		Fs: filesystem.NewBasicFileSystem(),
		Rw: rw,
	}
}

func NewReceiverTransfer(rw protoutils.ProtoMessageReadWriter, dest string) Transfer {
	return &receiver.ReceiveTransfer{
		Fs:   filesystem.NewBasicFileSystem(),
		Rw:   rw,
		Dest: dest,
	}
}
