package server

import (
	"iter"
	"net"
	"slices"
	"strings"
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

func (d exchangeDirection) String() string {
	switch d {
	case DirectionDial:
		return "dial"
	case DirectionListen:
		return "listen"
	default:
		return "unknown"
	}
}

// exchangeNodeMessage Used to exchange node information with other node
// If the connection is initiated by the other party (DirectionListen)
// We need to first read the node information sent by the other party.
// If the connection was initiated by us, we need to send local node info first.
func (s *RemoteXServer) exchangeNodeMessage(sc connection.Stream, direction exchangeDirection) (remote *node.Node, err error) {
	remote = new(node.Node)

	args := []*node.Node{remote, s.NodeService.GetLocal()}
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
	ipAddr := net.ParseIP(remote.Address.GetIpAddress())
	if ipAddr == nil || ipAddr.IsUnspecified() || ipAddr.IsLoopback() {
		addrSplit := strings.Split(sc.RemoteAddr().String(), ":")
		remote.Address.IpAddress = strings.Join(addrSplit[:len(addrSplit)-1], ":")
	}

	remote.IsLocal = false

	return remote, nil
}

func (s *RemoteXServer) connectionHandshake(conn connection.TlsConn) (*node.Node, error) {
	var (
		direction    exchangeDirection = DirectionDial
		streamGetter                   = conn.OpenStream
	)
	if conn.IsServer() {
		direction = DirectionListen
		streamGetter = conn.AcceptStream
	}

	stream, err := streamGetter()
	if err != nil {
		return nil, errors.Wrapf(err, "getStream: %v", direction)
	}

	defer func() {
		if direction == DirectionDial {
			stream.Close()
		}
	}()

	_ = stream.SetDeadline(time.Now().Add(time.Second * 2))

	remote, err := s.exchangeNodeMessage(stream, direction)
	if err != nil {
		return nil, errors.Wrap(err, "exchangeNode")
	}

	remote.ConnectionId = conn.GetConnectionId()
	remote.LastHeartbeat = time.Now().Unix()
	conn.SetNodeId(remote.NodeId)

	plog.Debugf("exchange remote node: %v", remote)

	return remote, s.AuthService.AuthConnection()
}
