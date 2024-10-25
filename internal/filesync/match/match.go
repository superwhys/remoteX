package match

import (
	"bytes"
	"context"
	"iter"

	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/hash"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/internal/filesystem"
)

// HashMatch is used to compare the HashHead transmitted from the client
// with the source file from the server
func HashMatch(ctx context.Context, head *pb.HashHead, srcFile *filesystem.File) (matchIter iter.Seq2[*pb.FileChunk, error], err error) {
	var (
		buf         []byte
		fileSize    int64
		offset      int64
		sum         uint32
		blockLength int64

		hashHits     int
		hashs        = head.GetHashs()
		hashMap      = GetHashMap(head)
		notMatchData = make([]byte, 0, head.GetBlockLength())
	)

	yieldNotMatchData := func(force bool, yield func(*pb.FileChunk, error) bool) {
		if force || len(notMatchData) >= int(head.GetBlockLength()) {
			yield(&pb.FileChunk{Data: notMatchData}, nil)
			notMatchData = notMatchData[:0]
		}
	}

	matchIter = func(yield func(*pb.FileChunk, error) bool) {
		fi, err := srcFile.Stat()
		if err != nil {
			yield(nil, err)
			return
		}

		fileSize = fi.Size()

		// read the first pb.FileChunk with offset: 0
		if buf, sum, blockLength, err = ReadFileBuf(srcFile, fileSize, offset, head); err != nil {
			yield(nil, err)
			return
		}

		defer yieldNotMatchData(true, yield)

	Outer:
		for {
			select {
			case <-ctx.Done():
				yield(nil, ctx.Err())
				return
			default:

			}
			blockIdx, match := hashMap[sum]
			if match {
				yieldNotMatchData(true, yield)
				for ; blockIdx < head.GetCheckSumCount(); blockIdx++ {
					select {
					case <-ctx.Done():
						yield(nil, ctx.Err())
						return
					default:
					}

					// check if the content of this block matches
					if hashs[blockIdx].GetLen() != blockLength || hashs[blockIdx].Sum1 != sum {
						break
					}

					sum2 := hash.CheckHashSum(buf)
					if !bytes.Equal(hashs[blockIdx].GetSum2(), sum2) {
						break
					}

					hashHits++
					if !yield(&pb.FileChunk{Hash: hashs[blockIdx]}, nil) {
						return
					}

					offset += blockLength
					if offset >= fileSize {
						break Outer
					}

					if buf, sum, blockLength, err = ReadFileBuf(srcFile, fileSize, offset, head); err != nil {
						yield(nil, err)
						return
					}
				}
			}

			notMatchData = append(notMatchData, buf[0])
			yieldNotMatchData(false, yield)

			offset++
			if offset >= fileSize {
				break
			}

			oldBuf := buf

			l := blockLength
			if remaining := fileSize - offset; remaining < l {
				l = remaining
			}
			buf, err = filesystem.BasicFs.ReadFileAtOffset(srcFile, offset, l)
			if err != nil {
				yield(nil, errors.Wrapf(err, "ReadFileAtOffset(%v-%v)", offset, l))
				return
			}

			sum = hash.RollingUpdate(sum, oldBuf[0], buf[l-1], uint32(blockLength))
		}
	}

	return matchIter, nil
}

func GetHashMap(head *pb.HashHead) map[uint32]int64 {
	hashMap := make(map[uint32]int64)
	for _, sum := range head.GetHashs() {
		hashMap[sum.GetSum1()] = sum.GetIndex()
	}
	return hashMap
}

func ReadFileBuf(f *filesystem.File, fileSize, offset int64, head *pb.HashHead) (buf []byte, sum uint32, blockLength int64, err error) {
	blockLength = head.GetBlockLength()
	if remaining := fileSize - offset; remaining < blockLength {
		blockLength = remaining
	}

	buf, err = filesystem.BasicFs.ReadFileAtOffset(f, offset, blockLength)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, "ReadFileAtOffset(%v-%v)", offset, blockLength)
	}
	sum = hash.CheckAdlerSum(buf)
	return buf, sum, blockLength, nil
}
