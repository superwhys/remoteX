package receiver

import (
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/remoteX/internal/filesync/hash"
	"github.com/superwhys/remoteX/internal/filesystem"
)

func TestCalcFileSubHash(t *testing.T) {
	tempFileCreate := func(size int) (*os.File, error) {
		tmpFile, err := os.CreateTemp("", "testfile-tmp-*.txt")
		if err != nil {
			return nil, err
		}
		
		content := make([]byte, size)
		for i := 0; i < size; i++ {
			content[i] = byte(i % 256)
		}
		
		_, err = tmpFile.Write(content)
		if err != nil {
			return nil, err
		}
		
		if err := tmpFile.Close(); err != nil {
			return nil, err
		}
		
		return tmpFile, nil
	}
	
	calc := func(f *os.File, size int) {
		defer os.Remove(f.Name())
		
		in, err := filesystem.BasicFs.OpenFile(f.Name())
		if !assert.Nil(t, err) {
			return
		}
		head := hash.CalcHashHead(int64(size))
		
		totalBlock := 0
		for _ = range hash.CalcFileSubHash(head, int64(size), in.File()) {
			totalBlock++
		}
		
		assert.Equal(t, head.CheckSumCount, int64(totalBlock))
	}
	
	t.Run("10MFile", func(t *testing.T) {
		size := 10 * 1024 * 1024 // 10MB
		f, err := tempFileCreate(size)
		if !assert.Nil(t, err) {
			return
		}
		calc(f, size)
	})
	
	t.Run("500MFile", func(t *testing.T) {
		size := 500 * 1024 * 1024 // 500MB
		f, err := tempFileCreate(size)
		if !assert.Nil(t, err) {
			return
		}
		calc(f, size)
	})
	
	t.Run("1024MFile", func(t *testing.T) {
		size := 1024 * 1024 * 1024 // 1G
		f, err := tempFileCreate(size)
		if !assert.Nil(t, err) {
			return
		}
		calc(f, size)
	})
}
