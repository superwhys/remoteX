package filesystem

import (
	"io/fs"
	"iter"
)

type FileSystem interface {
	List(path string) ([]*Entry, error)
	Walk(path string, filter WalkFilter) ([]*Entry, error)
	WalkIter(path string, filter WalkFilter) (iter.Seq[*Entry], error)
	OpenFile(path string) (File, error)
	CreateFile(path string) (File, error)
}

type WalkFilter func(path string, info fs.FileInfo) bool
