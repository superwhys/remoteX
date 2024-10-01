package filesystem

import (
	"context"
	
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/command/filesystem"
	"github.com/superwhys/remoteX/internal/fs"
	"github.com/superwhys/remoteX/pkg/errorutils"
)

type ServiceImpl struct {
	fs fs.FileSystem
}

func NewFilesystemService() filesystem.Service {
	return &ServiceImpl{
		fs: fs.NewBasicFileSystem(),
	}
}

func (s *ServiceImpl) ListDir(_ context.Context, args command.Args) (proto.Message, error) {
	path, exists := args["path"]
	if !exists {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Listdir), args)
	}
	
	entries, err := s.fs.List(path)
	if err != nil {
		return nil, err
	}
	
	return &fs.ListResp{
		Entries: entries,
	}, nil
}
