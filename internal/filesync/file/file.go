// File:		File.go
// Created by:	Hoven
// Created on:	2024-10-11
//
// This File is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package file

import (
	"os"
	
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesync/hash"
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

func ReadFileBuf(f *File, fileSize, offset int64, head *filesync.HashHead) (buf []byte, sum uint32, blockLength int64, err error) {
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
