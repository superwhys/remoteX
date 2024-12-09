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
				// if server node was down, it will try to reconnect to server again
				if !conn.IsServer() {
					s.connectionRedial(ctx, conn.GetNodeId(), conn.GetDialURL())
				}
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

	plog.Debugf("received command: %v", cmd)

	resp, err := s.CommandService.DoAcceptCommand(ctx, cmd, stream)
	if err != nil {
		return stream.WriteMessage(&command.Ret{ErrMsg: fmt.Sprintf("handle command failed: %v", err)})
	}

	return stream.WriteMessage(resp)
}

func (s *RemoteXServer) HandleLocalCommand(ctx context.Context, cmd *command.Command) (resp *command.Ret, err error) {
	return s.CommandService.DoLocalCommand(ctx, cmd)
}

func (s *RemoteXServer) SendCommandToRemote(ctx context.Context, nodeId common.NodeID, cmd *command.Command) (resp *command.Ret, err error) {
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

	if err := stream.WriteMessage(cmd); err != nil {
		return nil, errors.Wrap(err, "sendCommand")
	}

	resp = new(command.Ret)
	if err = stream.ReadMessage(resp); err != nil {
		return nil, errors.Wrap(err, "readRespMessage")
	}

	return resp, nil
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

	return s.CommandService.DoAcceptCommand(ctx, cmd, stream)
}

func (s *RemoteXServer) HandleCommandInBackground(ctx context.Context, nodeId common.NodeID, cmd *command.Command) (resp *command.Ret, err error) {
	remoteNode, err := s.NodeService.GetNode(nodeId)
	if err != nil {
		return nil, errors.Wrap(err, "getNode")
	}

	connId := remoteNode.GetConnectionId()

	conn, err := s.ConnService.GetConnection(connId)
	if err != nil {
		return nil, errors.Wrap(err, "getConnection")
	}

	return s.CommandService.DoOriginCommand(ctx, cmd, conn)
}
