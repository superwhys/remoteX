package command

import (
	"context"

	"github.com/superwhys/remoteX/domain/connection"
)

type Args map[string]Command_Arg

type RemoteOpt struct {
	Conn   connection.StreamConnection
	Stream connection.Stream
}

type Service interface {
	DoOriginCommand(ctx context.Context, cmd *Command, conn connection.StreamConnection) (*Ret, error)
	DoAcceptCommand(ctx context.Context, cmd *Command, stream connection.Stream) (*Ret, error)
	DoLocalCommand(ctx context.Context, cmd *Command) (*Ret, error)
}
