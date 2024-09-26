package server

import (
	"context"
	"fmt"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
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
				return errors.Wrap(err, "acceptStream")
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
	buffer := make([]byte, 1024)
	n, err := stream.Read(buffer)
	if err != nil {
		return errors.Wrap(err, "failed to read data")
	}
	
	plog.Infof("received data: %s", string(buffer[:n]))
	
	_, err = stream.Write([]byte(fmt.Sprintf("received: %s", string(buffer[:n]))))
	if err != nil {
		return errors.Wrap(err, "failed to write response")
	}
	
	return nil
}
