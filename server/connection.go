package server

import (
	"context"
	"iter"
	"net/url"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
	"github.com/superwhys/remoteX/pkg/common"
)

func (s *RemoteXServer) registerConnection(conn connection.TlsConn) {
	s.ConnService.RegisterConnection(conn)
}

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
				conn.Close()
				continue
			}
			
			remote, err := s.connectionHandshake(conn)
			if err != nil {
				plog.Errorf("listen connection handshake err: %v", err)
				conn.Close()
				continue
			}
			
			if conn.IsServer() {
				remote.Role = node.NodeConnectRoleServer
			} else {
				remote.Role = node.NodeConnectRoleClient
			}
			
			if !yield(remote, conn) {
				plog.Errorf("yield remoteNode and TlsConn error")
				conn.Close()
				break
			}
		}
	}
}

func (s *RemoteXServer) connectionRedial(nodeId common.NodeID, target *url.URL) {
	// TODO: retry while server node was down
}
