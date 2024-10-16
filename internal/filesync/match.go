package filesync

import (
	"bytes"
	"context"
	"iter"
	"os"
	"path/filepath"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
)

type chunk struct {
	// hash used to store block data that has already been matched
	hash *HashBuf
	// data used to store the origin byte data which not matched
	data []byte
}

func HashMatch(ctx context.Context, head *HashHead, strip string, fb *FileBase) (matchIter iter.Seq2[*chunk, error], err error) {
	var (
		f *os.File

		buf         []byte
		fileSize    int64
		offset      int64
		sum         uint32
		blockLength int64

		hashHits     int
		hashs        = head.GetHashs()
		hashMap      = getHashMap(head)
		notMatchData = make([]byte, 0, head.GetBlockLength())
	)

	yieldNotMatchData := func(force bool, yield func(*chunk, error) bool) {
		if force || len(notMatchData) >= int(head.GetBlockLength()) {
			yield(&chunk{data: notMatchData}, nil)
			notMatchData = notMatchData[:0]
		}
	}

	matchIter = func(yield func(*chunk, error) bool) {
		filePath := filepath.Join(strip, fb.GetPath())
		f, err = os.Open(filePath)
		if err != nil {
			yield(nil, err)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			yield(nil, err)
			return
		}

		fileSize = fi.Size()

		plog.Infof("HashSearch path=%s len(sums)=%d, srcFileSize=%d", filePath, len(head.GetHashs()), fileSize)

		// read the first chunk with offset: 0
		if buf, sum, blockLength, err = readFileBuf(f, fileSize, offset, head); err != nil {
			yield(nil, err)
			return
		}

		defer yieldNotMatchData(true, yield)

	Outer:
		for {
			blockIdx, match := hashMap[sum]
			if match {
				yieldNotMatchData(true, yield)
				for ; blockIdx < head.GetCheckSumCount(); blockIdx++ {
					// check if the content of this block matches
					if hashs[blockIdx].GetLen() != blockLength || hashs[blockIdx].Sum1 != sum {
						break
					}

					sum2 := CheckHashSum(buf)
					if !bytes.Equal(hashs[blockIdx].GetSum2(), sum2) {
						break
					}

					hashHits++
					if !yield(&chunk{hash: hashs[blockIdx]}, nil) {
						return
					}

					offset += blockLength
					if offset >= fileSize {
						break Outer
					}

					if buf, sum, blockLength, err = readFileBuf(f, fileSize, offset, head); err != nil {
						yield(nil, err)
						return
					}
				}
			}

			plog.Debugf("No match found for offset=%d", offset)
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
			buf, err = readFileAtOffset(f, offset, l)
			if err != nil {
				yield(nil, errors.Wrapf(err, "readFileAtOffset(%v-%v)", offset, l))
				return
			}

			sum = RollingUpdate(sum, oldBuf[0], buf[l-1], uint32(blockLength))
		}
	}

	return matchIter, nil
}
