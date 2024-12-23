// File:		proxy.go
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
	"net"
	"net/http"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/command/proxy"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"golang.org/x/sync/errgroup"
)

const (
	isServerKey  = "is_server"
	proxyAddrKey = "proxy_addr"
)

type ServiceImpl struct{}

func NewProxyService() proxy.Service {
	return &ServiceImpl{}
}

func (s *ServiceImpl) Name() string {
	return "proxy"
}

func (s *ServiceImpl) SupportedCommands() []command.CommandType {
	return []command.CommandType{
		command.Proxy,
	}
}

func (s *ServiceImpl) Invoke(ctx context.Context, cmd *command.Command, cmdCtx *command.CommandContext) (proto.Message, error) {
	switch cmd.Type {
	case command.Proxy:
		err := s.Proxy(ctx, cmd.GetArgs(), cmdCtx.Remote, cmdCtx.ExchangeWriter)
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, errorutils.ErrCommandTypeNotSupport(cmd.GetType().String())
	}
}

func (s *ServiceImpl) Proxy(ctx context.Context, args command.Args, stream connection.Stream, exchangeWriter io.Writer) error {
	isServerArg, exists := args.GetArg("is_server")
	if !exists {
		return errorutils.ErrCommandMissingArguments(command.Proxy.String(), isServerKey)
	}

	switch isServerArg.GetBoolValue() {
	case true:
		return s.proxyServer(ctx, args, stream, exchangeWriter)
	case false:
		return s.proxyClient(ctx, args, stream, exchangeWriter)
	}

	return nil
}

func (s *ServiceImpl) proxyServer(_ context.Context, args command.Args, stream connection.Stream, exchangeWriter io.Writer) error {
	args[isServerKey] = command.BoolArg(false)
	cmd := &command.Command{
		Type: command.Proxy,
		Args: args,
	}

	err := stream.WriteMessage(cmd)
	if err != nil {
		return errorutils.WrapRemoteXError(err, "writeProxyServerCmd")
	}

	hijacker, ok := exchangeWriter.(http.Hijacker)
	if !ok {
		return errorutils.WrapRemoteXError(nil, "hijacking not supported")
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		return errorutils.WrapRemoteXError(err, "getHijackConn")
	}

	defer conn.Close()
	if _, err = io.Copy(conn, stream); err != nil {
		return errorutils.WrapRemoteXError(err, "copyConnToStream")
	}
	return nil
}

func (s *ServiceImpl) proxyClient(ctx context.Context, args command.Args, stream connection.Stream, exchangeWriter io.Writer) error {
	proxyAddrArg, exists := args.GetArg(proxyAddrKey)
	if !exists {
		return errorutils.ErrCommandMissingArguments(command.Proxy.String(), proxyAddrKey)
	}

	conn, err := net.Dial("tcp", proxyAddrArg.GetStrValue())
	if err != nil {
		return errorutils.WrapRemoteXError(err, "dialRemoteAddr")
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-ctx.Done()
		_ = conn.Close()
		plog.Debugc(ctx, "forwardService: close connections")
		return stream.Close()
	})
	eg.Go(func() error {
		defer plog.Debugc(ctx, "forwardService:copy to remote over")
		_, err := io.Copy(conn, stream)
		if err != nil {
			return errorutils.WrapRemoteXError(err, "copyStreamToConn")
		}
		return nil
	})
	eg.Go(func() error {
		defer plog.Debugc(ctx, "forwardService:copy to tcp connection over")
		_, err := io.Copy(stream, conn)
		if err != nil {
			return errorutils.WrapRemoteXError(err, "copyConnToStream")
		}
		return nil
	})

	return eg.Wait()
}
