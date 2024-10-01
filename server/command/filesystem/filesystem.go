package filesystem

import (
	"github.com/superwhys/remoteX/domain/command/filesystem"
	"github.com/superwhys/remoteX/internal/fs"
)

type ServiceImpl struct {
	fs fs.FileSystem
}

func NewFilesystemService() filesystem.Service {
	return &ServiceImpl{
		fs: fs.NewBasicFileSystem(),
	}
}

func (s *ServiceImpl) ListDir(path string) (*fs.ListResp, error) {
	entries, err := s.fs.List(path)
	if err != nil {
		return nil, err
	}
	
	return &fs.ListResp{
		Entries: entries,
	}, nil
}
