package hash

import (
	"crypto/sha256"
	"io"
	"iter"
	"math"
	"os"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/internal/filesync"
)

const blockSize = 64

func CalcHashHead(contentLen int64) *filesync.HashHead {
	blockLength := int64(math.Sqrt(float64(contentLen)))
	if blockLength < blockSize {
		blockLength = blockSize
	}
	
	const checkSumLength = 32
	
	return &filesync.HashHead{
		CheckSumCount:   (contentLen + (int64(blockLength) - 1)) / int64(blockLength),
		RemainderLength: contentLen % int64(blockLength),
		BlockLength:     blockLength,
		CheckSumLength:  checkSumLength,
	}
}

func CheckHashSum(buf []byte) []byte {
	hash := sha256.New()
	
	hash.Write(buf)
	return hash.Sum(nil)
}

type FileChunk struct {
	chunkIdx int64
	offset   int64
	b        []byte
}

func CalcFileSubHash(head *filesync.HashHead, fileLen int64, in *os.File) iter.Seq[*filesync.HashBuf] {
	return func(yield func(*filesync.HashBuf) bool) {
		for fc := range splitFileChunk(head.BlockLength, fileLen, head.CheckSumCount, in) {
			sum1 := CheckAdlerSum(fc.b)
			sum2 := CheckHashSum(fc.b)
			hb := &filesync.HashBuf{
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

func splitFileChunk(blockLen, fileLen, chunkCnt int64, in *os.File) iter.Seq[*FileChunk] {
	return func(yield func(*FileChunk) bool) {
		buf := make([]byte, int(blockLen))
		remaining := fileLen
		var offset int64
		
		for i := int64(0); i < chunkCnt; i++ {
			n1 := min(blockLen, remaining)
			b := buf[:n1]
			
			if _, err := io.ReadFull(in, b); err != nil {
				plog.Errorf("read File: %v block error. idx: %v, offset: %v err: %v", in.Name(), i, offset, err)
				break
			}
			
			fc := &FileChunk{
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
