package filesystem

import (
	"crypto/md5"
	"io"
	"os"

	"github.com/pkg/errors"
)

var _ File = (*BaseFile)(nil)

type File interface {
	File() *os.File
	Name() string
	Stat() (os.FileInfo, error)
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Seek(offset int64, whence int) (n int64, err error)
	Update(from File)
	CurrentSeek() (int64, error)
	MD5() ([]byte, error)
	ReadFileAtOffset(offset int64, length int64) ([]byte, error)
	Close() error
}

type BaseFile struct {
	file *os.File
}

func (f *BaseFile) File() *os.File {
	return f.file
}

func (f *BaseFile) Stat() (os.FileInfo, error) {
	return f.file.Stat()
}

func (f *BaseFile) Name() string {
	return f.file.Name()
}

func (f *BaseFile) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *BaseFile) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func (f *BaseFile) Seek(offset int64, whence int) (n int64, err error) {
	return f.file.Seek(offset, whence)
}

func (f *BaseFile) Update(from File) {
	f.file = from.File()
}

func (f *BaseFile) CurrentSeek() (int64, error) {
	return f.file.Seek(0, io.SeekCurrent)
}

func (f *BaseFile) MD5() ([]byte, error) {
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

func (f *BaseFile) ReadFileAtOffset(offset int64, length int64) ([]byte, error) {
	buffer := make([]byte, length)

	_, err := f.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	n, err := f.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

func (f *BaseFile) Close() error {
	return f.file.Close()
}
