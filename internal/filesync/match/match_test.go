// File:		match_test.go
// Created by:	Hoven
// Created on:	2024-10-12
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package match

import (
	"context"
	"crypto/md5"
	"os"
	"slices"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/remoteX/internal/filesync/file"
	"github.com/superwhys/remoteX/internal/filesync/hash"
)

func tempFileCreate(size int) (string, error) {
	tmpFile, err := os.CreateTemp("", "testfile-tmp-*.txt")
	if err != nil {
		return "", err
	}
	
	content := make([]byte, size)
	for i := 0; i < size; i++ {
		content[i] = byte(i % 256)
	}
	
	_, err = tmpFile.Write(content)
	if err != nil {
		return "", err
	}
	
	if err := tmpFile.Close(); err != nil {
		return "", err
	}
	
	return tmpFile.Name(), nil
}

func TestFileHashMatch(t *testing.T) {
	target := "./content/test.txt"
	source := "./content/src.txt"
	targetFile, err := file.OpenFile(target)
	if !assert.Nil(t, err) {
		return
	}
	defer targetFile.Close()
	src, err := file.OpenFile(source)
	if !assert.Nil(t, err) {
		return
	}
	defer src.Close()
	
	fi, err := targetFile.Stat()
	if !assert.Nil(t, err) {
		return
	}
	size := fi.Size()
	
	head := hash.CalcHashHead(size)
	head.Hashs = slices.Collect(hash.CalcFileSubHash(head, size, targetFile.File()))
	
	// In actual, matchIter data should be transmitted back to the client
	matchIter, err := HashMatch(context.Background(), head, src)
	if !assert.Nil(t, err) {
		return
	}
	
	// In actual, The client received the chunks transmitted by the server's matchIter
	// and concatenated them into a complete targetFile based on the chunks
	md5Hash := md5.New()
	err = SyncFileToWriter(matchIter, targetFile, md5Hash)
	if !assert.Nil(t, err) {
		return
	}
	
	srcMd5, err := src.MD5()
	if !assert.Nil(t, err) {
		return
	}
	
	finalHash := md5Hash.Sum(nil)
	assert.Equal(t, srcMd5, finalHash)
}
