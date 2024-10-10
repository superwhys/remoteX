package server

import (
	"iter"
	"slices"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
)

type exchangeDirection int

const (
	DirectionDial = iota
	DirectionListen
)

// exchangeNodeMessage Used to exchange node information with other node
// If the connection is initiated by the other party (DirectionListen)
// We need to first read the node information sent by the other party.
// If the connection was initiated by us, we need to send local node info first.
func (s *RemoteXServer) exchangeNodeMessage(sc connection.Stream, direction exchangeDirection) (remote *node.Node, err error) {
	remote = new(node.Node)

	args := []*node.Node{remote, s.nodeService.GetLocal()}
	fns := []func(message proto.Message) error{
		sc.ReadMessage,
		sc.WriteMessage,
	}

	var iterFn iter.Seq2[int, func(message proto.Message) error]
	switch direction {
	case DirectionDial:
		iterFn = slices.Backward(fns)
	case DirectionListen:
		iterFn = slices.All(fns)
	default:
		return nil, errors.New("invalid direction")
	}

	for idx, fn := range iterFn {
		arg := args[idx]
		if err := fn(arg); err != nil {
			return nil, err
		}
	}

	remote.IsLocal = false

	return remote, nil
}

func (s *RemoteXServer) connectionHandshake(conn connection.TlsConn) (*node.Node, error) {
	var (
		direction    exchangeDirection
		streamGetter = conn.OpenStream
	)
	if conn.IsServer() {
		direction = DirectionListen
		streamGetter = conn.AcceptStream
	}

	stream, err := streamGetter()
	if err != nil {
		return nil, errors.Wrap(err, "acceptStream")
	}

	defer stream.Close()

	_ = stream.SetDeadline(time.Now().Add(time.Second * 2))

	remote, err := s.exchangeNodeMessage(stream, direction)
	if err != nil {
		return nil, errors.Wrap(err, "exchangeNode with client")
	}

	remote.ConnectionId = conn.GetConnectionId()
	remote.LastHeartbeat = time.Now().Unix()
	conn.SetNodeId(remote.NodeId)

	plog.Debugf("exchange remote node: %v", remote)

	return remote, s.authService.AuthConnection()
}
