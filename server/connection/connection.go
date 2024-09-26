package connection

import (
	"context"
	"crypto/tls"
	"iter"
	"time"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
)

type ConnectionAppService struct {
	service     connection.Service
	local       *node.Node
	connections chan connection.TlsConn
}

func NewConnectionAppService(local *node.Node, tlsConf *tls.Config) *ConnectionAppService {
	return &ConnectionAppService{
		local:       local,
		service:     connection.NewConnectionService(local.URL(), tlsConf),
		connections: make(chan connection.TlsConn),
	}
}

func (c *ConnectionAppService) CreateListener(ctx context.Context) error {
	return c.service.CreateListener(ctx, c.connections)
}

func (c *ConnectionAppService) ExchangeNodeMessage(tc connection.TlsConn) (remote *node.Node, err error) {
	remote = new(node.Node)
	if err = tc.ReadMessage(remote); err != nil {
		return nil, err
	}
	
	if err = tc.WriteMessage(c.local); err != nil {
		return nil, err
	}
	
	return remote, nil
}

func (c *ConnectionAppService) CloseConn(conn connection.TlsConn) {
	c.service.CloseConnection(conn.GetConnectionId())
}

func (c *ConnectionAppService) HandleConnection(ctx context.Context) iter.Seq[connection.TlsConn] {
	return func(yield func(connection.TlsConn) bool) {
		for {
			var conn connection.TlsConn
			select {
			case <-ctx.Done():
				break
			case conn = <-c.connections:
			}
			
			if err := c.service.CheckConnection(conn); err != nil {
				plog.Errorf("check connection err: %v", err)
				c.CloseConn(conn)
				continue
			}
			
			_ = conn.SetDeadline(time.Now().Add(20 * time.Second))
			
			yield(conn)
		}
	}
}
