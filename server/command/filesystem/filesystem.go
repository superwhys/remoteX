package filesystem

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/internal/filesystem"
	"github.com/superwhys/remoteX/pkg/errorutils"

	fsDomain "github.com/superwhys/remoteX/domain/command/filesystem"
)

type ServiceImpl struct {
	fs filesystem.FileSystem
}

func NewFilesystemService() fsDomain.Service {
	return &ServiceImpl{
		fs: filesystem.NewBasicFileSystem(),
	}
}

func (s *ServiceImpl) Invoke(ctx context.Context, cmd *command.Command, opt *command.RemoteOpt) (proto.Message, error) {
	switch cmd.Type {
	case command.Listdir:
		return s.ListDir(ctx, cmd.GetArgs(), nil)
	default:
		return nil, errorutils.ErrNotSupportCommandType(int32(cmd.Type))
	}
}

func (s *ServiceImpl) ListDir(_ context.Context, args command.Args, _ *command.RemoteOpt) (proto.Message, error) {
	path, exists := args["path"]
	if !exists {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Listdir), args)
	}

	entries, err := s.fs.List(path.GetStrValue())
	if err != nil {
		return nil, err
	}

	return &filesystem.ListResp{
		Entries: entries,
	}, nil
}
