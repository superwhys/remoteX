// File:		service.go
// Created by:	Hoven
// Created on:	2024-09-30
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package filesystem

import (
	"github.com/superwhys/remoteX/internal/fs"
)

type Service interface {
	ListDir(path string) (*fs.ListResp, error)
}

type ServiceImpl struct {
	fs fs.FileSystem
}

func NewFilesystemService() Service {
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
