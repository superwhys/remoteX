// File:		service.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package screenshot

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type Service interface {
	Screenshot(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error)
}
