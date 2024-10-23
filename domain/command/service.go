package command

import (
	"context"

	"github.com/superwhys/remoteX/domain/connection"
)

type Args map[string]Command_Arg

type Service interface {
	DoCommand(ctx context.Context, cmd *Command, stream connection.Stream) (*Ret, error)
}
