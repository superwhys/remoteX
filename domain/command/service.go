package command

import (
	"context"
)

type Args map[string]string

type Service interface {
	DoCommand(ctx context.Context, cmd *Command) (*Ret, error)
}
