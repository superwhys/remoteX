package filesync

import (
	"crypto/sha256"
	"math"
)

const blockSize = 64

func calcHashHead(contentLen int64) *HashHead {
	blockLength := int64(math.Sqrt(float64(contentLen)))
	if blockLength < blockSize {
		blockLength = blockSize
	}

	const checkSumLength = 32

	return &HashHead{
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

func getHashMap(head *HashHead) map[uint32]int64 {
	hashMap := make(map[uint32]int64)
	for _, sum := range head.GetHashs() {
		hashMap[sum.GetSum1()] = sum.GetIndex()
	}
	return hashMap
}

type hashTarget struct {
	index int64
	sum   *HashBuf
}

func getHashTargets(head *HashHead) []*hashTarget {
	targets := make([]*hashTarget, len(head.GetHashs()))
	for _, sum := range head.GetHashs() {
		targets[sum.GetIndex()] = &hashTarget{
			index: sum.GetIndex(),
			sum:   sum,
		}
	}
	return targets
}
