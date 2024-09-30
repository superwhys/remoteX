package server

import (
	"context"
	"iter"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/auth"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
	"github.com/superwhys/remoteX/pkg/connection/dialer"
	"github.com/superwhys/remoteX/pkg/connection/listener"
	"github.com/superwhys/remoteX/pkg/limiter"
	"github.com/superwhys/remoteX/pkg/svcutils"
	"github.com/thejerf/suture/v4"
	"golang.org/x/sync/errgroup"
)

type RemoteXServer struct {
	*suture.Supervisor

	opt               *Option
	nodeService       node.Service
	connService       connection.Service
	authService       auth.Service
	commandService    command.Service
	limiter           *limiter.Limiter
	heartbeatInterval time.Duration
	// connections It is a channel used for transmitting successfully established connections
	// (including listening connections and self created connections)
	connections chan connection.TlsConn
}

func NewRemoteXServer(opt *Option) *RemoteXServer {
	local := opt.Local
	transConf := opt.Conf.TransConf
	server := &RemoteXServer{
		opt:        opt,
		Supervisor: suture.NewSimple("RemoteX.Service"),

		nodeService:       node.NewNodeService(local),
		authService:       auth.NewSimpleAuthService(),
		connService:       connection.NewConnectionService(local.URL(), opt.TlsConfig),
		commandService:    command.NewCommandService(),
		heartbeatInterval: time.Second * time.Duration(opt.Conf.HeartbeatInterval),
		limiter:           limiter.NewLimiter(local.NodeId, transConf.MaxRecvKbps, transConf.MaxSendKbps),
		connections:       make(chan connection.TlsConn),
	}

	server.registerNode(opt.Local)

	listener.InitListener()
	dialer.InitDialer()

	go server.StartDialer(context.Background())

	server.Add(svcutils.AsService(server.StartListener, "startListener"))
	server.Add(svcutils.AsService(server.HandleConnection, "handleConnection"))

	return server
}

func (s *RemoteXServer) StartDialer(ctx context.Context) error {
	for _, client := range s.opt.Conf.DialClients {
		target := client.URL()

		conn, err := s.connService.EstablishConnection(ctx, target)
		if err != nil {
			return errors.Wrap(err, "failed to establish connection")
		}

		s.connections <- conn
	}

	return nil
}

func (s *RemoteXServer) StartListener(ctx context.Context) error {
	return s.connService.CreateListener(ctx, s.connections)
}

func (s *RemoteXServer) HandleConnection(ctx context.Context) error {
	for remote, conn := range s.connectionPrepareIter(ctx) {
		go func(remote *node.Node, conn connection.TlsConn) {
			if err := s.registerNode(remote); err != nil {
				plog.Errorf("register node: %v: %v", remote, err)
				return
			}
			s.registerConnection(conn)
			plog.Debugf("register connection: %v. NodeId: %v", conn.GetConnectionId(), conn.GetNodeId())

			s.background(ctx, conn, conn.IsServer())
		}(remote, conn)
	}
	return nil
}

func (s *RemoteXServer) connectionPrepareIter(ctx context.Context) iter.Seq2[*node.Node, connection.TlsConn] {
	return func(yield func(*node.Node, connection.TlsConn) bool) {
		for {
			var conn connection.TlsConn
			select {
			case <-ctx.Done():
				break
			case conn = <-s.connections:
			}

			if err := s.connService.CheckConnection(conn); err != nil {
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

			if !yield(remote, conn) {
				plog.Errorf("yield remoteNode and TlsConn error")
				conn.Close()
				break
			}
		}
	}
}

func (s *RemoteXServer) background(ctx context.Context, conn connection.TlsConn, isServer bool) {
	eg, ctx := errgroup.WithContext(ctx)

	hbStartNotify := make(chan struct{})

	eg.Go(func() error {
		return s.schedulerHeartbeat(ctx, conn, hbStartNotify, isServer)
	})

	// must be called after heartbeat is start
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-hbStartNotify:
			return s.schedulerCommand(ctx, conn)
		}
	})

	if err := eg.Wait(); err != nil {
		plog.Errorf("failed to run background connection: %v", err)
		s.nodeService.UpdateNodeStatus(conn.GetNodeId(), node.NodeStatusOffline)
		s.CloseConn(conn)
	}
}

func (s *RemoteXServer) registerNode(n *node.Node) error {
	return s.nodeService.RegisterNode(n)
}

func (s *RemoteXServer) registerConnection(conn connection.TlsConn) {
	s.connService.RegisterConnection(conn)
}

func (s *RemoteXServer) CloseConn(conn connection.TlsConn) {
	s.connService.CloseConnection(conn.GetConnectionId())
}
