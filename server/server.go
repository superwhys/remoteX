package server

import (
	"context"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/pkg/svcutils"
	"github.com/superwhys/remoteX/server/auth"
	"github.com/superwhys/remoteX/server/connection"
	"github.com/superwhys/remoteX/server/node"
	"github.com/thejerf/suture/v4"
)

type RemoteXServer struct {
	*suture.Supervisor
	
	opt     *Option
	connApp *connection.ConnectionAppService
	authApp *auth.AuthAppService
	nodeApp *node.NodeAppService
}

func NewRemoteXServer(opt *Option) *RemoteXServer {
	server := &RemoteXServer{
		opt:        opt,
		Supervisor: suture.NewSimple("RemoteX.Service"),
		connApp:    connection.NewConnectionAppService(opt.Local, opt.TlsConfig),
		authApp:    auth.NewAuthAppService(opt.Local),
	}
	
	server.Add(svcutils.AsService(server.StartListener, "startListener"))
	server.Add(svcutils.AsService(server.HandleConnection, "handleConnection"))
	
	return server
}

func (s *RemoteXServer) StartListener(ctx context.Context) error {
	return s.connApp.CreateListener(ctx)
}

func (s *RemoteXServer) HandleConnection(ctx context.Context) error {
	for c := range s.connApp.HandleConnection(ctx) {
		conn := c
		go func() {
			// 1. handshake
			remote, err := s.connApp.ExchangeNodeMessage(conn)
			if err != nil {
				plog.Errorf("exchange node: %v message: %v", conn.RemoteAddr(), err)
				s.connApp.CloseConn(conn)
				return
			}
			
			plog.Debugf("exchange remote node: %v", remote)
			// 2. auth
			if err := s.authApp.AuthRemoteConn(remote, conn); err != nil {
				plog.Errorf("auth remote node: %v: %v", remote, err)
				s.connApp.CloseConn(conn)
				return
			}
			// 3. register node
			remote.ConnectionId = conn.GetConnectionId()
			if err := s.nodeApp.RegisterNode(remote); err != nil {
				plog.Errorf("register node: %v: %v", remote, err)
				s.connApp.CloseConn(conn)
				return
			}
		}()
	}
	
	return nil
}
