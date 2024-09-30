package server

import (
	"context"
	"net"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/counter"
)

func (s *RemoteXServer) schedulerCommand(ctx context.Context, conn connection.TlsConn) error {
	for {
		select {
		case <-ctx.Done():
			plog.Errorf("context done: %v", ctx.Err())
			return ctx.Err()
		default:
			stream, err := conn.AcceptStream()
			if err != nil {
				opErr := new(net.OpError)
				if errors.As(err, &opErr) && !opErr.Timeout() {
					return errors.Wrap(err, "acceptStream")
				}

				continue
			}

			// pack limiter and counter
			limitStream := connection.PackLimiterStream(stream, s.limiter)
			counterStream := connection.PackCounterConnection(
				stream,
				&counter.CountingReader{Reader: limitStream},
				&counter.CountingWriter{Writer: limitStream},
			)

			if err := s.handleStream(ctx, counterStream); err != nil {
				return errors.Wrap(err, "handleStream")
			}
		}
	}
}

func (s *RemoteXServer) handleStream(_ context.Context, stream connection.Stream) error {
	command := &command.Command{}
	if err := stream.ReadMessage(command); err != nil {
		return errors.Wrap(err, "readMessage")
	}

	plog.Infof("received command: %v", command)
	return nil
}
