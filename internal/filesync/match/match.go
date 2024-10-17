package match

import (
	"bytes"
	"context"
	"iter"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesync/file"
	"github.com/superwhys/remoteX/internal/filesync/hash"
)

type chunk struct {
	// hash used to store block data that has already been matched
	hash *filesync.HashBuf
	// data used to store the origin byte data which not matched
	data []byte
}

// HashMatch is used to compare the HashHead transmitted from the client
// with the source file from the server
func HashMatch(ctx context.Context, head *filesync.HashHead, srcFile *file.File) (matchIter iter.Seq2[*chunk, error], err error) {
	var (
		buf         []byte
		fileSize    int64
		offset      int64
		sum         uint32
		blockLength int64
		
		hashHits     int
		hashs        = head.GetHashs()
		hashMap      = filesync.GetHashMap(head)
		notMatchData = make([]byte, 0, head.GetBlockLength())
	)
	
	yieldNotMatchData := func(force bool, yield func(*chunk, error) bool) {
		if force || len(notMatchData) >= int(head.GetBlockLength()) {
			yield(&chunk{data: notMatchData}, nil)
			notMatchData = notMatchData[:0]
		}
	}
	
	matchIter = func(yield func(*chunk, error) bool) {
		fi, err := srcFile.Stat()
		if err != nil {
			yield(nil, err)
			return
		}
		
		fileSize = fi.Size()
		
		plog.Infof("HashSearch path=%s len(sums)=%d, srcFileSize=%d", srcFile.Name(), len(head.GetHashs()), fileSize)
		
		// read the first chunk with offset: 0
		if buf, sum, blockLength, err = file.ReadFileBuf(srcFile, fileSize, offset, head); err != nil {
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
					
					sum2 := hash.CheckHashSum(buf)
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
					
					if buf, sum, blockLength, err = file.ReadFileBuf(srcFile, fileSize, offset, head); err != nil {
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
			buf, err = file.ReadFileAtOffset(srcFile, offset, l)
			if err != nil {
				yield(nil, errors.Wrapf(err, "ReadFileAtOffset(%v-%v)", offset, l))
				return
			}
			
			sum = hash.RollingUpdate(sum, oldBuf[0], buf[l-1], uint32(blockLength))
		}
	}
	
	return matchIter, nil
}
