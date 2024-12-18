package server

import (
	"context"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
	"github.com/superwhys/remoteX/pkg/errorutils"
)

func (s *RemoteXServer) schedulerHeartbeat(ctx context.Context, conn connection.TlsConn, hbStartNotify chan struct{}) error {
	var (
		streamGetter func() (connection.Stream, error)
		handler      func(stream connection.Stream) error
	)

	switch conn.IsServer() {
	case true:
		streamGetter = conn.AcceptStream
		handler = s.receiveHeartbeat
	case false:
		streamGetter = conn.OpenStream
		handler = s.sendHeartbeat
	}

	hbStream, err := streamGetter()
	if err != nil {
		return errorutils.ErrGetHeartbeatStream(err)
	}
	defer hbStream.Close()

	// close the channel to notify commandScheduler
	close(hbStartNotify)

	ticket := time.NewTicker(s.heartbeatInterval)
	defer ticket.Stop()

	for {
		if err := hbStream.SetDeadline(time.Now().Add(s.heartbeatInterval)); err != nil {
			return errorutils.WrapRemoteXError(err, "failed to set write deadline")
		}

		if err := handler(hbStream); err != nil {
			return errorutils.WrapRemoteXError(err, "failed to handle heartbeat")
		}

		conn.UpdateLastHeartbeat()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticket.C:
		}
	}
}

func (s *RemoteXServer) sendHeartbeat(stream connection.Stream) error {
	cn, err := s.NodeService.RefreshCurrentNode()
	if err != nil {
		return errors.Wrap(err, "failed to refresh current node")
	}
	cn.IsLocal = false

	if err := stream.WriteMessage(cn); err != nil {
		return errorutils.WrapRemoteXError(err, "send heartbeat")
	}

	plog.Debugf("send heartbeat to %v", stream.RemoteAddr())

	return nil
}

func (s *RemoteXServer) receiveHeartbeat(stream connection.Stream) error {
	rn := new(node.Node)
	if err := stream.ReadMessage(rn); err != nil {
		return errorutils.WrapRemoteXError(err, "read heartbeat")
	}

	if err := s.NodeService.UpdateHeartbeat(stream.GetNodeId()); err != nil {
		return errors.Wrap(err, "failed to update heartbeat")
	}

	plog.Debugf("receive heartbeat from node: %v remoteAddr: %v", rn.NodeId, stream.RemoteAddr())
	return nil
}
