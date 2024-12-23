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
		err := s.Proxy(ctx, cmd.GetArgs(), cmdCtx.Remote, cmdCtx.ProxyReq, cmdCtx.ExchangeWriter)
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, errorutils.ErrCommandTypeNotSupport(cmd.GetType().String())
	}
}

func (s *ServiceImpl) Proxy(ctx context.Context, args command.Args, stream connection.Stream, req *http.Request, exchangeWriter io.Writer) error {
	isServerArg, exists := args.GetArg("is_server")
	if !exists {
		return errorutils.ErrCommandMissingArguments(command.Proxy.String(), isServerKey)
	}

	switch isServerArg.GetBoolValue() {
	case true:
		return s.proxyServer(ctx, args, stream, req, exchangeWriter)
	case false:
		return s.proxyClient(ctx, args, stream)
	}

	return nil
}

func (s *ServiceImpl) proxyServer(ctx context.Context, args command.Args, stream connection.Stream, req *http.Request, exchangeWriter io.Writer) error {
	args[isServerKey] = command.BoolArg(false)
	cmd := &command.Command{
		Type: command.Proxy,
		Args: args,
	}

	err := stream.WriteMessage(cmd)
	if err != nil {
		return errorutils.WrapRemoteXError(err, "writeProxyServerCmd")
	}

	err = req.Write(stream)
	if err != nil {
		return errorutils.WrapRemoteXError(err, "writeProxyServerReq")
	}

	hijacker, ok := exchangeWriter.(http.Hijacker)
	if !ok {
		return errorutils.WrapRemoteXError(nil, "hijacking not supported")
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		return errorutils.WrapRemoteXError(err, "getHijackConn")
	}
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()
		_ = conn.Close()
		plog.Debugc(ctx, "proxyHijack: close connections")
		return stream.Close()
	})

	eg.Go(func() error {
		defer plog.Debugc(ctx, "proxyServer: copy stream to conn")
		_, err := io.Copy(conn, stream)
		if err != nil {
			return errorutils.WrapRemoteXError(err, "copyStreamToConn")
		}
		return nil
	})
	eg.Go(func() error {
		defer plog.Debugc(ctx, "proxyServer: copy conn to stream")
		_, err := io.Copy(stream, conn)
		if err != nil {
			return errorutils.WrapRemoteXError(err, "copyConnToStream")
		}
		return nil
	})

	return eg.Wait()
}

func (s *ServiceImpl) proxyClient(ctx context.Context, args command.Args, stream connection.Stream) error {
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
		plog.Debugc(ctx, "proxyClient: close connections")
		return stream.Close()
	})
	eg.Go(func() error {
		defer plog.Debugc(ctx, "proxyClient: copy stream to conn")
		_, err := io.Copy(conn, stream)
		if err != nil {
			return errorutils.WrapRemoteXError(err, "copyStreamToConn")
		}
		return nil
	})
	eg.Go(func() error {
		defer plog.Debugc(ctx, "proxyClient: copy conn to stream")
		_, err := io.Copy(stream, conn)
		if err != nil {
			return errorutils.WrapRemoteXError(err, "copyConnToStream")
		}
		return nil
	})

	return eg.Wait()
}
