package server

import (
	"context"
	"fmt"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"golang.org/x/sync/errgroup"
)

func (s *RemoteXServer) schedulerCommand(ctx context.Context, conn connection.TlsConn) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			stream, err := conn.AcceptStream()
			if errorutils.IsRemoteDead(err) {
				plog.Errorf("remote(%v) was down", conn.RemoteAddr())
				// if server node was down, it will try to reconnect to server again
				if !conn.IsServer() {
					s.connectionRedial(conn.GetNodeId(), conn.GetDialURL())
				}
				return err
			} else if err != nil {
				continue
			}

			go func(stream connection.Stream) {
				// pack tracker, limiter and counter
				stream = connection.PackStream(stream, s.packOpts)

				if err := s.handleCommand(ctx, stream); err != nil {
					plog.Errorf("handle command error: %v", err)
				}
			}(stream)
		}
	}
}

func (s *RemoteXServer) handleCommand(ctx context.Context, stream connection.Stream) error {
	cmd := &command.Command{}
	if err := stream.ReadMessage(cmd); err != nil {
		return errors.Wrap(err, "readMessage")
	}

	plog.Debugf("received command: %v", cmd)

	resp, err := s.commandService.DoCommand(ctx, cmd, stream)
	if err != nil {
		if resp == nil {
			resp = &command.Ret{}
		}

		if resp.ErrMsg == "" {
			resp.ErrMsg = fmt.Sprintf("handle command failed: %v", err)
		}
	}

	if err := stream.WriteMessage(resp); err != nil {
		return errors.Wrap(err, "writeMessage")
	}

	return nil
}

type callbackFn func(ctx context.Context, stream connection.Stream) error

func (s *RemoteXServer) handleRemoteCommand(ctx context.Context, nodeId common.NodeID, cmd *command.Command, callback callbackFn) (resp *command.Ret, err error) {
	remoteNode, err := s.nodeService.GetNode(nodeId)
	if err != nil {
		return nil, errors.Wrap(err, "getNode")
	}

	connId := remoteNode.GetConnectionId()

	conn, err := s.connService.GetConnection(connId)
	if err != nil {
		return nil, errors.Wrap(err, "getConnection")
	}

	stream, err := conn.OpenStream()
	if err != nil {
		return nil, errors.Wrap(err, "openStream")
	}
	defer stream.Close()

	stream = connection.PackStream(stream, s.packOpts)
	return s.handleRemoteStream(ctx, stream, cmd, callback)
}

func (s *RemoteXServer) handleRemoteStream(ctx context.Context, stream connection.Stream, cmd *command.Command, callback callbackFn) (resp *command.Ret, err error) {

	eg, ctx := errgroup.WithContext(ctx)

	if callback != nil {
		eg.Go(func() error {
			return callback(ctx, stream)
		})
	}

	eg.Go(func() error {
		return stream.WriteMessage(cmd)
	})

	if err := eg.Wait(); err != nil {
		return nil, errors.Wrap(err, "do remote command")
	}

	resp = new(command.Ret)
	if err = stream.ReadMessage(resp); err != nil {
		return nil, errors.Wrap(err, "readRespMessage")
	}

	return resp, nil
}
