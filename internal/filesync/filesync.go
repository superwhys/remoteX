package filesync

import (
	"path/filepath"
)

func (f *FileBase) Name() string {
	return filepath.Clean(string(f.Path))
}

func GetHashMap(head *HashHead) map[uint32]int64 {
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

func GetHashTargets(head *HashHead) []*hashTarget {
	targets := make([]*hashTarget, len(head.GetHashs()))
	for _, sum := range head.GetHashs() {
		targets[sum.GetIndex()] = &hashTarget{
			index: sum.GetIndex(),
			sum:   sum,
		}
	}
	return targets
}
