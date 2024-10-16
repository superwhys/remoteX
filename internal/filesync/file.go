// File:		file.go
// Created by:	Hoven
// Created on:	2024-10-11
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package filesync

import (
	io "io"
	"iter"
	"os"
	"path/filepath"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
)

func (f *FileBase) Name() string {
	return filepath.Clean(string(f.Path))
}

func readFileAtOffset(file *os.File, offset int64, length int64) ([]byte, error) {
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

func readFileBuf(f *os.File, fileSize, offset int64, head *HashHead) (buf []byte, sum uint32, blockLength int64, err error) {
	blockLength = head.GetBlockLength()
	if remaining := fileSize - offset; remaining < blockLength {
		blockLength = remaining
	}

	buf, err = readFileAtOffset(f, offset, blockLength)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, "readFileAtOffset(%v-%v)", offset, blockLength)
	}
	sum = CheckAdlerSum(buf)
	return buf, sum, blockLength, nil
}

type fileChunk struct {
	chunkIdx int64
	offset   int64
	b        []byte
}

func splitFileChunk(blockLen, fileLen, chunkCnt int64, in *os.File) iter.Seq[*fileChunk] {
	return func(yield func(*fileChunk) bool) {
		buf := make([]byte, int(blockLen))
		remaining := fileLen
		var offset int64

		for i := int64(0); i < chunkCnt; i++ {
			n1 := min(blockLen, remaining)
			b := buf[:n1]

			if _, err := io.ReadFull(in, b); err != nil {
				plog.Errorf("read file: %v block error. idx: %v, offset: %v err: %v", in.Name(), i, offset, err)
				break
			}

			fc := &fileChunk{
				chunkIdx: i,
				offset:   offset,
				b:        b,
			}

			if !yield(fc) {
				break
			}

			remaining -= n1
			offset += n1
		}
	}
}

func calcFileSubHash(head *HashHead, fileLen int64, in *os.File) iter.Seq[*HashBuf] {
	return func(yield func(*HashBuf) bool) {
		for fc := range splitFileChunk(head.BlockLength, fileLen, head.CheckSumCount, in) {
			sum1 := CheckAdlerSum(fc.b)
			sum2 := CheckHashSum(fc.b)
			hb := &HashBuf{
				Index:  fc.chunkIdx,
				Offset: fc.offset,
				Sum1:   sum1,
				Sum2:   sum2,
			}
			if fc.chunkIdx == head.CheckSumCount-1 && head.RemainderLength != 0 {
				hb.Len = int64(head.RemainderLength)
			} else {
				hb.Len = int64(head.BlockLength)
			}

			if !yield(hb) {
				break
			}
		}
	}
}

func SyncFile(matchIter iter.Seq2[*chunk, error], target *os.File, writer io.Writer) error {
	for matchChunk, err := range matchIter {
		if err != nil {
			return errors.Wrap(err, "iter file match")
		}

		var data []byte
		if matchChunk.hash != nil {
			offset := matchChunk.hash.GetOffset()
			data, err = readFileAtOffset(target, offset, matchChunk.hash.GetLen())
			if err != nil {
				return errors.Wrap(err, "read file at offset")
			}
		} else {
			data = matchChunk.data
		}

		_, err := writer.Write(data)
		if err != nil {
			return errors.Wrap(err, "write to Writer")
		}
	}

	return nil
}
