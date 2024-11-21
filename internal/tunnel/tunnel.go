package tunnel

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
	"golang.org/x/sync/errgroup"
)

type Tunnel struct {
	TunnelKey  string
	LocalAddr  string
	RemoteAddr string
	listener   net.Listener
	conn       connection.StreamConnection
}

func (t *Tunnel) Close() error {
	if t.listener != nil {
		t.listener.Close()
	}
	if t.conn != nil {
		t.conn.Close()
	}
	return nil
}

type TunnelManager struct {
	tunnels map[string]*Tunnel
	mu      sync.RWMutex
}

func NewTunnelManager() *TunnelManager {
	return &TunnelManager{
		tunnels: make(map[string]*Tunnel),
	}
}

func (m *TunnelManager) CreateTunnel(ctx context.Context, localAddr, remoteAddr string, conn connection.StreamConnection) (*Tunnel, error) {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %v", localAddr, err)
	}

	tunnel := &Tunnel{
		TunnelKey:  uuid.New().String(),
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
		listener:   listener,
		conn:       conn,
	}

	m.mu.Lock()
	m.tunnels[tunnel.TunnelKey] = tunnel
	m.mu.Unlock()

	go m.handleTunnel(ctx, tunnel)

	return tunnel, nil
}

func (m *TunnelManager) handleTunnel(ctx context.Context, tunnel *Tunnel) {
	defer m.CloseTunnel(tunnel.TunnelKey)

	for {
		select {
		case <-ctx.Done():
			plog.Debugc(ctx, "tunnel %s context canceled, closing tunnel", tunnel.TunnelKey)
			return
		default:
			if err := tunnel.listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second)); err != nil {
				plog.Errorc(ctx, "Failed to set accept deadline: %v", err)
				return
			}

			localConn, err := tunnel.listener.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue
				}
				plog.Errorc(ctx, "Failed to accept connection: %v", err)
				return
			}

			go func(ctx context.Context, conn net.Conn) {
				defer conn.Close()

				stream, err := m.remoteEstablish(tunnel)
				if err != nil {
					plog.Errorc(ctx, "Failed to establish connection: %v", err)
					return
				}
				defer stream.Close()
				plog.Debugc(ctx, "establish remote connection success, start exchange data")

				if err := m.exchange(ctx, stream, conn); err != nil {
					plog.Errorc(ctx, "Failed to exchange data: %v", err)
				}
			}(ctx, localConn)
		}
	}
}

func (m *TunnelManager) remoteEstablish(tunnel *Tunnel) (connection.Stream, error) {
	tc := &command.TunnelConnect{
		TunnelKey: tunnel.TunnelKey,
		Addr:      tunnel.RemoteAddr,
	}

	stream, err := tunnel.conn.OpenStream()
	if err != nil {
		return nil, err
	}

	if err := stream.WriteMessage(tc.ToCommand(command.Forwardreceive)); err != nil {
		return nil, errors.Wrap(err, "send forwardreceive command")
	}

	resp := &command.TunnelConnectResp{}
	if err := stream.ReadMessage(resp); err != nil {
		return nil, errors.Wrap(err, "read remote forwardreceive command resp")
	}

	if !resp.Success {
		return nil, errors.New("remote establish not success")
	}

	return stream, nil
}

func (m *TunnelManager) ReceiveTunnel(ctx context.Context, tunnelKey, targetAddr string, stream connection.Stream) error {
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		return stream.WriteMessage(&command.TunnelConnectResp{
			TunnelKey: tunnelKey,
			Success:   false,
			Error:     err.Error(),
		})
	}
	defer targetConn.Close()

	if err := stream.WriteMessage(&command.TunnelConnectResp{
		TunnelKey: tunnelKey,
		Success:   true,
	}); err != nil {
		return err
	}

	return m.exchange(ctx, stream, targetConn)
}

func (m *TunnelManager) ListTunnels() []*Tunnel {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tunnels := make([]*Tunnel, 0, len(m.tunnels))
	for _, tunnel := range m.tunnels {
		tunnels = append(tunnels, tunnel)
	}
	return tunnels
}

func (m *TunnelManager) CloseTunnel(tunnelKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if tunnel, exists := m.tunnels[tunnelKey]; exists {
		tunnel.Close()
		delete(m.tunnels, tunnelKey)
	}
}

func (m *TunnelManager) copyWithContext(ctx context.Context, src io.Reader, dst io.Writer) error {
	buf := make([]byte, 1024*16)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			l, err := src.Read(buf)
			if l > 0 {
				if _, we := dst.Write(buf[:l]); err != nil {
					return we
				}
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
		}
	}
}

func (m *TunnelManager) exchange(ctx context.Context, dst, src io.ReadWriter) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return m.copyWithContext(ctx, dst, src)
	})

	eg.Go(func() error {
		return m.copyWithContext(ctx, src, dst)
	})

	return eg.Wait()
}
