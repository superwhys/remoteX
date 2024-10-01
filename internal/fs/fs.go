package fs

import (
	"os"
	"path/filepath"
	"time"
)

type FileSystem interface {
	List(path string) ([]*Entry, error)
}

var _ FileSystem = (*BasicFileSystem)(nil)

type BasicFileSystem struct{}

func NewBasicFileSystem() *BasicFileSystem {
	return &BasicFileSystem{}
}

func (f *BasicFileSystem) List(path string) (entries []*Entry, err error) {
	if !filepath.IsAbs(path) {
		return nil, ErrDirPahtNotAbs
	}

	if !PathExists(path) {
		return nil, ErrDirPathNotFound
	}

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, de := range dirEntries {
		info, err := de.Info()
		if err != nil {
			return nil, err
		}

		entry := &Entry{
			Name:         de.Name(),
			Path:         filepath.Join(path, de.Name()),
			Size:         info.Size(),
			Type:         GetEntryType(de.IsDir()),
			CreatedTime:  info.ModTime().Format(time.DateTime),
			ModifiedTime: info.ModTime().Format(time.DateTime),
			AccessedTime: info.ModTime().Format(time.DateTime),
			Owner:        "",
			Permissions:  "",
		}
		entry.Owner, entry.Permissions, _ = getFileInfo(path)
		entries = append(entries, entry)
	}

	return
}
