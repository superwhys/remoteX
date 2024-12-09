package server

import (
	"context"
	"iter"
	"math"
	"net/url"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
	"github.com/superwhys/remoteX/pkg/common"
)

func (s *RemoteXServer) CloseConnection(conn connection.TlsConn) {
	s.ConnService.CloseConnection(conn.GetConnectionId())
}

func (s *RemoteXServer) connectionHandshakeIter(ctx context.Context) iter.Seq2[*node.Node, connection.TlsConn] {
	return func(yield func(*node.Node, connection.TlsConn) bool) {
		for {
			var conn connection.TlsConn
			select {
			case <-ctx.Done():
				break
			case conn = <-s.connections:
			}

			if err := s.ConnService.CheckConnection(conn); err != nil {
				plog.Errorf("check connection err: %v", err)
				s.CloseConnection(conn)
				continue
			}

			remote, err := s.connectionHandshake(conn)
			if err != nil {
				plog.Errorf("listen connection handshake err: %v", err)
				s.CloseConnection(conn)
				if !conn.IsServer() {
					s.connectionRedial(ctx, remote.NodeId, remote.URL())
				}
				continue
			}

			if conn.IsServer() {
				remote.Role = node.NodeConnectRoleServer
			} else {
				remote.Role = node.NodeConnectRoleClient
			}

			if !yield(remote, conn) {
				plog.Errorf("yield remoteNode and TlsConn error")
				s.CloseConnection(conn)
				break
			}
		}
	}
}

// connectionRedial reconnects connections that have been established but were disconnected during operation
func (s *RemoteXServer) connectionRedial(ctx context.Context, nodeId common.NodeID, target *url.URL) {
	s.connectionRedialByTask(ctx, connection.NewDialTask(target, nodeId, time.Second*2, s.maxDialDelay, true))
}

// connectionRedialByTask reconnects the DialTask that is currently failed and sets the connection count to +1,
// which is used for reconnection when the connection task fails
func (s *RemoteXServer) connectionRedialByTask(ctx context.Context, task *connection.DialTask) {
	task.DialCnt++
	task.IsRedial = true

	var delay time.Duration
	attempt := task.DialCnt
	if attempt >= task.Threshold {
		delay = task.MaxDelay
	} else {
		delay = time.Duration(math.Min(
			float64(task.MaxDelay),
			float64(task.InitDelay)*math.Pow(2, float64(attempt)),
		))
	}
	plog.Errorc(ctx, "Attempt %d failed, retrying in %v...", attempt+1, delay)

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
			s.dialTasks <- task
		}
	}()
}
