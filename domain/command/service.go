package command

import (
	"context"
)

type Service interface {
	DoCommand(ctx context.Context, cmd *Command) (*Ret, error)
}
