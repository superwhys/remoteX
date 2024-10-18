package filesystem

import (
	"crypto/md5"
	"io"
	"os"
	
	"github.com/pkg/errors"
)

type File struct {
	file *os.File
}

func (f *File) File() *os.File {
	return f.file
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.file.Stat()
}

func (f *File) Name() string {
	return f.file.Name()
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *File) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func (f *File) Seek(offset int64, whence int) (n int64, err error) {
	return f.file.Seek(offset, whence)
}

func (f *File) Update(from *File) {
	f.file = from.file
}

func (f *File) CurrentSeek() (int64, error) {
	return f.file.Seek(0, io.SeekCurrent)
}

func (f *File) MD5() ([]byte, error) {
	m := md5.New()
	
	currentSeek, err := f.CurrentSeek()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current seek")
	}
	if _, err := f.file.Seek(0, io.SeekStart); err != nil {
		return nil, errors.Wrap(err, "failed to seek file")
	}
	if _, err := io.Copy(m, f.file); err != nil {
		return nil, errors.Wrap(err, "failed to copy file")
	}
	
	if _, err := f.file.Seek(currentSeek, io.SeekStart); err != nil {
		return nil, errors.Wrap(err, "failed to restore seek")
	}
	
	return m.Sum(nil), nil
}

func (f *File) Close() error {
	return f.file.Close()
}
