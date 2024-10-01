// File:		service.go
// Created by:	Hoven
// Created on:	2024-09-30
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package filesystem

import (
	"context"
	
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
)

type Service interface {
	ListDir(ctx context.Context, args command.Args) (proto.Message, error)
}
