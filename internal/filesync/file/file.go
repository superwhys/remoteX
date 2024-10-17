package file

import (
	"crypto/md5"
	"io"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/hash"
	"github.com/superwhys/remoteX/internal/filesync/pb"
)

type File struct {
	file *os.File
}

func OpenFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	
	return &File{
		file: f,
	}, nil
}

func CreateFile(path string) (*File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	
	return &File{
		file: f,
	}, nil
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

func ReadFileAtOffset(file *File, offset int64, length int64) ([]byte, error) {
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

func ReadFileBuf(f *File, fileSize, offset int64, head *pb.HashHead) (buf []byte, sum uint32, blockLength int64, err error) {
	blockLength = head.GetBlockLength()
	if remaining := fileSize - offset; remaining < blockLength {
		blockLength = remaining
	}
	
	buf, err = ReadFileAtOffset(f, offset, blockLength)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, "ReadFileAtOffset(%v-%v)", offset, blockLength)
	}
	sum = hash.CheckAdlerSum(buf)
	return buf, sum, blockLength, nil
}

func GetFileList(strip, root string) iter.Seq[*pb.FileBase] {
	ch := make(chan *pb.FileBase)
	go func() {
		filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				plog.Errorf("filepath walk error: %v", err)
				return err
			}
			if path == root && info.IsDir() {
				return nil
			}
			
			name := strings.TrimPrefix(path, strip)
			if name == root {
				name = "."
			}
			
			size := info.Size()
			
			f := &pb.FileBase{
				Name:    name,
				Regular: info.Mode().IsRegular(),
				Wpath:   name,
				Size:    size,
			}
			ch <- f
			
			return nil
		})
		
		close(ch)
	}()
	
	return func(yield func(*pb.FileBase) bool) {
		for f := range ch {
			if !yield(f) {
				break
			}
		}
	}
}
