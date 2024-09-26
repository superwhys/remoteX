package server

import (
	"context"
	"time"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
)

func (s *RemoteXServer) schedulerHeartbeat(ctx context.Context, conn connection.TlsConn, isServer bool) error {
	var (
		streamGetter func(conn connection.TlsConn) (connection.Stream, error)
		handler      func(stream connection.Stream) error
	)
	
	switch isServer {
	case true:
		streamGetter = func(conn connection.TlsConn) (connection.Stream, error) {
			return conn.AcceptStream()
		}
		handler = s.receiveHeartbeat
	case false:
		streamGetter = func(conn connection.TlsConn) (connection.Stream, error) {
			return conn.OpenStream()
		}
		handler = s.sendHeartbeat
	default:
		return errors.New("invalid connection type")
	}
	
	hbStream, err := streamGetter(conn)
	if err != nil {
		return errors.Wrap(err, "failed to get heartbeat stream")
	}
	defer hbStream.Close()
	
	ticket := time.NewTicker(s.heartbeatInterval)
	defer ticket.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticket.C:
		}
		
		if err := hbStream.SetDeadline(time.Now().Add(s.heartbeatInterval)); err != nil {
			return errors.Wrap(err, "failed to set write deadline")
		}
		
		if err := handler(hbStream); err != nil {
			return errors.Wrap(err, "failed to handle heartbeat")
		}
		
		conn.UpdateLastHeartbeat()
	}
}

func (s *RemoteXServer) sendHeartbeat(stream connection.Stream) error {
	cn, err := s.nodeService.RefreshCurrentNode()
	if err != nil {
		return errors.Wrap(err, "failed to refresh current node")
	}
	
	if err := stream.WriteMessage(cn); err != nil {
		return errors.Wrap(err, "failed to write heartbeat message")
	}
	
	plog.Debugf("send heartbeat to %v", stream.RemoteAddr())
	
	return nil
}

func (s *RemoteXServer) receiveHeartbeat(stream connection.Stream) error {
	rn := new(node.Node)
	if err := stream.ReadMessage(rn); err != nil {
		return errors.Wrap(err, "failed to read heartbeat message")
	}
	
	if err := s.nodeService.UpdateHeartbeat(stream.GetNodeId()); err != nil {
		return errors.Wrap(err, "failed to update heartbeat")
	}
	
	plog.Debugf("receive heartbeat from node: %v remoteAddr: %v", rn.NodeId, stream.RemoteAddr())
	return nil
}
