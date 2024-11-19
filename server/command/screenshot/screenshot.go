// File:		screenshot.go
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
	"github.com/superwhys/remoteX/domain/command/screenshot"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

var _ screenshot.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
}

func NewScreenshotService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s *ServiceImpl) Screenshot(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error) {
	return nil, nil
}
