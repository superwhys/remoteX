package filesync

import (
	"testing"

	"github.com/go-puzzles/puzzles/plog"
)

func TestCalcHashHead(t *testing.T) {
	t.Run("Divisible", func(t *testing.T) {
		size := 100 * 1024 * 1024 // 100 MB
		t.Logf("file size: %v", size)
		head := calcHashHead(int64(size))
		t.Log(plog.Jsonify(head))
	})

	t.Run("InDivisible", func(t *testing.T) {
		size := 100 * 1024 * 1024 // 100 MB
		size = size + 1111
		t.Logf("file size: %v", size)
		head := calcHashHead(int64(size))
		t.Log(plog.Jsonify(head))
	})

	t.Run("test10M", func(t *testing.T) {
		size := 10 * 1024 * 1024 // 10 MB
		t.Logf("file size: %v", size)
		head := calcHashHead(int64(size))
		t.Log(plog.Jsonify(head))
	})
}
