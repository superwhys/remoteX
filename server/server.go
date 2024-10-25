package server

import (
	"context"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/auth"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
	"github.com/superwhys/remoteX/internal/connection/dialer"
	"github.com/superwhys/remoteX/internal/connection/listener"
	"github.com/superwhys/remoteX/pkg/limiter"
	"github.com/superwhys/remoteX/pkg/svcutils"
	"github.com/superwhys/remoteX/pkg/tracker"
	"github.com/thejerf/suture/v4"
	"golang.org/x/sync/errgroup"

	authSrv "github.com/superwhys/remoteX/server/auth"
	commandSrv "github.com/superwhys/remoteX/server/command"
	connSrv "github.com/superwhys/remoteX/server/connection"
	nodeSrv "github.com/superwhys/remoteX/server/node"
)

type RemoteXServer struct {
	*suture.Supervisor

	NodeService    node.Service
	ConnService    connection.Service
	AuthService    auth.Service
	CommandService command.Service

	opt               *Option
	packOpts          *connection.PackOpts
	heartbeatInterval time.Duration

	dialTasks chan *connection.DialTask
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

		NodeService:       nodeSrv.NewNodeService(local),
		AuthService:       authSrv.NewSimpleAuthService(),
		ConnService:       connSrv.NewConnectionService(local.URL(), opt.TlsConfig),
		CommandService:    commandSrv.NewCommandService(),
		heartbeatInterval: time.Second * time.Duration(opt.Conf.HeartbeatInterval),
		packOpts: &connection.PackOpts{
			Limiter:        limiter.NewLimiter(local.NodeId, transConf.MaxRecvKbps, transConf.MaxSendKbps),
			TrackerManager: tracker.CreateTrackerManager(),
		},
		dialTasks:   make(chan *connection.DialTask),
		connections: make(chan connection.TlsConn),
	}

	listener.InitListener()
	dialer.InitDialer()

	server.Add(svcutils.AsService(server.startDialer, "startDialer"))
	server.Add(svcutils.AsService(server.startListener, "startListener"))
	server.Add(svcutils.AsService(server.handleConnection, "handleConnection"))

	return server
}

func (s *RemoteXServer) startDialer(ctx context.Context) error {
	go func() {
		for _, client := range s.opt.Conf.DialClients {
			s.dialTasks <- &connection.DialTask{
				Target: client.URL(),
			}
		}
	}()

	for {
		var task *connection.DialTask
		select {
		case <-ctx.Done():
			close(s.dialTasks)
			return ctx.Err()
		case task = <-s.dialTasks:
		}

		// if dial task is redial task
		// update the node status to Connecting
		if task.IsRedial && task.NodeId != "" {
			_ = s.NodeService.UpdateNodeStatus(task.NodeId, node.NodeStatusConnecting)
		}

		conn, err := s.ConnService.EstablishConnection(ctx, task.Target)
		if err != nil {
			if !task.IsRedial {
				return errors.Wrap(err, "failed to establish connection")
			}
			_ = s.NodeService.UpdateNodeStatus(task.NodeId, node.NodeStatusOffline)
			s.connectionRedial(task.NodeId, task.Target)
			continue
		}

		s.connections <- conn
	}
}

func (s *RemoteXServer) startListener(ctx context.Context) error {
	return s.ConnService.CreateListener(ctx, s.connections)
}

func (s *RemoteXServer) handleConnection(ctx context.Context) error {
	for remote, conn := range s.connectionHandshakeIter(ctx) {
		// check whether the remote node is existing
		// if exists, just update the nodeInfo
		var err error
		n, _ := s.NodeService.GetNode(remote.NodeId)
		if n != nil {
			err = s.NodeService.UpdateNode(remote)
		} else {
			err = s.NodeService.RegisterNode(remote)
		}
		if err != nil {
			plog.Errorf("update/register remote node: %v error: %v", remote, err)
			conn.Close()
			continue
		}
		s.ConnService.RegisterConnection(conn)
		plog.Infof("register connection: %v. NodeId: %v", conn.GetConnectionId(), conn.GetNodeId())

		go s.background(ctx, conn)
	}
	return nil
}

func (s *RemoteXServer) background(ctx context.Context, conn connection.TlsConn) {
	eg, ctx := errgroup.WithContext(ctx)

	hbStartNotify := make(chan struct{})

	eg.Go(func() error {
		return s.schedulerHeartbeat(ctx, conn, hbStartNotify)
	})

	// must be called after heartbeat is start
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-hbStartNotify:
			plog.Debugf("conn{%s} start command handler", conn.GetConnectionId())
			return s.schedulerCommand(ctx, conn)
		}
	})

	if err := eg.Wait(); err != nil {
		plog.Errorf("failed to run background connection: %v", err)
		_ = s.NodeService.UpdateNodeStatus(conn.GetNodeId(), node.NodeStatusOffline)
		s.CloseConnection(conn)
	}
}
