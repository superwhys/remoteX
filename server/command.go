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
				return err
			} else if err != nil {
				continue
			}

			go func(stream connection.Stream) {
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

	plog.Debugc(ctx, "received command: %v", cmd)
	defer plog.Debugc(ctx, "received command: %v done", cmd)

	resp, err := s.CommandService.Execute(ctx, cmd, &command.CommandContext{Remote: stream})
	if err != nil {
		return stream.WriteMessage(&command.Ret{ErrMsg: fmt.Sprintf("handle command failed: %v", err)})
	}

	return stream.WriteMessage(resp)
}

func (s *RemoteXServer) HandleLocalCommand(ctx context.Context, cmd *command.Command) (resp *command.Ret, err error) {
	return s.CommandService.Execute(ctx, cmd, &command.CommandContext{})
}

func (s *RemoteXServer) HandleCommandWithRemote(ctx context.Context, nodeId common.NodeID, cmd *command.Command) (resp *command.Ret, err error) {
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

	return s.CommandService.Execute(ctx, cmd, &command.CommandContext{
		IsRemote: true,
		Remote:   stream,
	})
}

func (s *RemoteXServer) HandleCommandWithRawRemote(ctx context.Context, nodeId common.NodeID, cmd *command.Command) (resp *command.Ret, err error) {
	remoteNode, err := s.NodeService.GetNode(nodeId)
	if err != nil {
		return nil, errors.Wrap(err, "getNode")
	}

	connId := remoteNode.GetConnectionId()

	conn, err := s.ConnService.GetConnection(connId)
	if err != nil {
		return nil, errors.Wrap(err, "getConnection")
	}

	return s.CommandService.Execute(ctx, cmd, &command.CommandContext{
		IsRemote:  true,
		RawRemote: conn,
	})
}
