package filesystem

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type Service interface {
	ListDir(ctx context.Context, args command.Args, rw protoutils.ProtoMessageReadWriter) (proto.Message, error)
}
