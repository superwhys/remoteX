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
				if err := s.receiveAndHandleCommand(ctx, stream); err != nil {
					plog.Errorf("handle command error: %v", err)
				}
			}(stream)
		}
	}
}

func (s *RemoteXServer) receiveAndHandleCommand(ctx context.Context, stream connection.Stream) error {
	cmd := &command.Command{}
	if err := stream.ReadMessage(cmd); err != nil {
		stream.Close()
		return errors.Wrap(err, "receiveCommand")
	}

	plog.Debugf("received command: %v", cmd)

	resp, err := s.HandleCommand(ctx, cmd, stream)
	if err != nil {
		return stream.WriteMessage(&command.Ret{ErrMsg: fmt.Sprintf("handle command failed: %v", err)})
	}

	return stream.WriteMessage(resp)
}

func (s *RemoteXServer) HandleCommand(ctx context.Context, cmd *command.Command, stream connection.Stream) (*command.Ret, error) {
	return s.CommandService.DoCommand(ctx, cmd, stream)
}

type callbackFn func(ctx context.Context, stream connection.Stream) error

func (s *RemoteXServer) HandleRemoteCommand(ctx context.Context, nodeId common.NodeID, cmd *command.Command, callback callbackFn) (resp *command.Ret, err error) {
	remoteNode, err := s.NodeService.GetNode(nodeId)
	if err != nil {
		return nil, errors.Wrap(err, "getNode")
	}

	connId := remoteNode.GetConnectionId()

	conn, err := s.ConnService.GetConnection(connId)
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

	if err := stream.WriteMessage(cmd); err != nil {
		return nil, errors.Wrap(err, "sendCommand")
	}

	if callback != nil {
		eg.Go(func() error {
			return callback(ctx, stream)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.Wrap(err, "do remote command")
	}

	resp = new(command.Ret)
	if err = stream.ReadMessage(resp); err != nil {
		return nil, errors.Wrap(err, "readRespMessage")
	}

	return resp, nil
}
