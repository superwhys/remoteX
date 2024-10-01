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
