package server

import (
	"context"
	"fmt"
	"net"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/common"
)

func (s *RemoteXServer) schedulerCommand(ctx context.Context, conn connection.TlsConn) error {
	for {
		select {
		case <-ctx.Done():
			plog.Errorf("context done: %v", ctx.Err())
			return ctx.Err()
		default:
			stream, err := conn.AcceptStream()
			if err != nil {
				opErr := new(net.OpError)
				if errors.As(err, &opErr) && !opErr.Timeout() {
					return errors.Wrap(err, "acceptStream")
				}

				continue
			}

			go func(stream connection.Stream) {
				// pack limiter and counter
				limitStream := connection.PackLimiterStream(stream, s.limiter)
				counterStream := connection.PackCounterStream(limitStream)

				if err := s.handleCommand(ctx, counterStream); err != nil {
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

	resp, err := s.commandService.DoCommand(ctx, cmd)
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

func (s *RemoteXServer) handleRemoteCommand(ctx context.Context, nodeId common.NodeID, cmd *command.Command) (resp *command.Ret, err error) {
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

	if err = stream.WriteMessage(cmd); err != nil {
		return nil, errors.Wrap(err, "writeCommandMessage")
	}

	resp = new(command.Ret)
	if err = stream.ReadMessage(resp); err != nil {
		return nil, errors.Wrap(err, "readRespMessage")
	}

	return resp, nil
}
