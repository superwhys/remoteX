// File:		receiver_test.go
// Created by:	Hoven
// Created on:	2024-10-11
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package receiver

import (
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/remoteX/internal/filesync/file"
	"github.com/superwhys/remoteX/internal/filesync/hash"
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
	
	t.Run("smallFile", func(t *testing.T) {
		size := 10 * 1024 * 1024 // 10MB
		f, err := tempFileCreate(size)
		if !assert.Nil(t, err) {
			return
		}
		defer os.Remove(f.Name())
		
		in, err := file.OpenFile(f.Name())
		if !assert.Nil(t, err) {
			return
		}
		head := hash.CalcHashHead(int64(size))
		
		totalBlock := 0
		for _ = range hash.CalcFileSubHash(head, int64(size), in.File()) {
			totalBlock++
		}
		
		assert.Equal(t, head.CheckSumCount, int64(totalBlock))
	})
}
