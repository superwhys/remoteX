package filesystem

import (
	"context"
	
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
)

type Service interface {
	ListDir(ctx context.Context, args command.Args) (proto.Message, error)
}
