package filesystem

import (
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	BasicFs *BasicFileSystem
	once    sync.Once
)

var _ FileSystem = (*BasicFileSystem)(nil)

type FileSystem interface {
	List(path string) ([]*Entry, error)
	Walk(path string, filter WalkFilter) ([]*Entry, error)
	WalkIter(path string, filter WalkFilter) (iter.Seq[*Entry], error)
	OpenFile(path string) (*File, error)
	CreateFile(path string) (*File, error)
	ReadFileAtOffset(file *File, offset int64, length int64) ([]byte, error)
}

type WalkFilter func(path string, info fs.FileInfo) bool

type BasicFileSystem struct{}

func NewBasicFileSystem() *BasicFileSystem {
	once.Do(func() {
		BasicFs = &BasicFileSystem{}
	})

	return BasicFs
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

		path := filepath.Join(path, de.Name())

		isDir, err := f.checkWheatherDir(path, info)
		if err != nil {
			return nil, err
		}

		entry := &Entry{
			Name:         de.Name(),
			Path:         path,
			Size:         info.Size(),
			Type:         GetEntryType(isDir),
			CreatedTime:  info.ModTime().Format(time.DateTime),
			ModifiedTime: info.ModTime().Format(time.DateTime),
			AccessedTime: info.ModTime().Format(time.DateTime),
			Regular:      info.Mode().IsRegular(),
			Owner:        "",
			Permissions:  "",
		}
		entry.Owner, entry.Permissions, _ = getFileInfo(path)
		entries = append(entries, entry)
	}

	return
}

func (f *BasicFileSystem) Walk(path string, filter WalkFilter) (entries []*Entry, err error) {
	if !filepath.IsAbs(path) {
		return nil, ErrDirPahtNotAbs
	}

	if !PathExists(path) {
		return nil, ErrDirPathNotFound
	}

	return
}

func (f *BasicFileSystem) checkWheatherDir(path string, info fs.FileInfo) (bool, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	if fileInfo.Mode()&os.ModeSymlink == 0 {
		return info.IsDir(), nil
	}

	targetPath, err := os.Readlink(path)
	if err != nil {
		return false, err
	}

	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(filepath.Dir(path), targetPath)
	}

	targetInfo, err := os.Stat(targetPath)
	if err != nil {
		return false, err
	}

	return targetInfo.IsDir(), nil
}

func (f *BasicFileSystem) WalkIter(root string, filter WalkFilter) (iter.Seq[*Entry], error) {
	if !filepath.IsAbs(root) {
		return nil, ErrDirPahtNotAbs
	}

	if !PathExists(root) {
		return nil, ErrDirPathNotFound
	}

	ch := make(chan *Entry)
	go func() {
		strip := filepath.Dir(filepath.Clean(root)) + "/"
		if strings.HasSuffix(root, "/") {
			strip = filepath.Clean(root) + "/"
		}

		filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filter != nil && !filter(path, info) {
				return nil
			}

			isDir, err := f.checkWheatherDir(path, info)
			if err != nil {
				return err
			}

			entry := &Entry{
				Name:         info.Name(),
				Type:         GetEntryType(isDir),
				Size:         info.Size(),
				Path:         path,
				Wpath:        strings.TrimPrefix(path, strip),
				CreatedTime:  info.ModTime().Format(time.DateTime),
				ModifiedTime: info.ModTime().Format(time.DateTime),
				AccessedTime: info.ModTime().Format(time.DateTime),
				Regular:      info.Mode().IsRegular(),
			}

			entry.Owner, entry.Permissions, _ = getFileInfo(path)
			ch <- entry
			return nil
		})
		close(ch)
	}()

	return func(yield func(*Entry) bool) {
		for e := range ch {
			if !yield(e) {
				break
			}
		}
	}, nil
}

func (f *BasicFileSystem) OpenFile(path string) (*File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	// file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &File{
		file: file,
	}, nil
}

func (f *BasicFileSystem) CreateFile(path string) (*File, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &File{
		file: file,
	}, nil
}

func (f *BasicFileSystem) ReadFileAtOffset(file *File, offset int64, length int64) ([]byte, error) {
	buffer := make([]byte, length)

	_, err := file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	n, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}
