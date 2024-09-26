package server

import (
	"time"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
)

func (s *RemoteXServer) connectionHandshake(conn connection.TlsConn) (*node.Node, error) {
	var (
		direction    exchangeDirection
		streamGetter = conn.OpenStream
	)
	if conn.IsServer() {
		direction = DirectionListen
		streamGetter = conn.AcceptStream
	}
	_ = conn.SetDeadline(time.Now().Add(time.Second))
	
	stream, err := streamGetter()
	if err != nil {
		return nil, errors.Wrap(err, "acceptStream")
	}
	
	defer stream.Close()
	
	remote, err := s.exchangeNodeMessage(stream, direction)
	if err != nil {
		return nil, errors.Wrap(err, "exchangeNode with client")
	}
	
	remote.ConnectionId = conn.GetConnectionId()
	conn.SetNodeId(remote.NodeId)
	
	plog.Debugf("exchange remote node: %v", remote)
	
	return remote, s.authService.AuthConnection()
}
