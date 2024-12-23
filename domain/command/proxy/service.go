// File:		service.go
// Created by:	Hoven
// Created on:	2024-12-21
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package proxy

import (
	"context"
	"io"
	"net/http"

	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
)

type Service interface {
	command.CommandProvider
	Proxy(ctx context.Context, args command.Args, stream connection.Stream, req *http.Request, w io.Writer) error
}
