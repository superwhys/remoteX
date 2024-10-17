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
	"os"
	"slices"
	"testing"
	
	"github.com/stretchr/testify/assert"
	file2 "github.com/superwhys/remoteX/internal/filesync/file"
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
	file, err := file2.OpenFile(target)
	if !assert.Nil(t, err) {
		return
	}
	defer file.Close()
	src, err := file2.OpenFile(source)
	if !assert.Nil(t, err) {
		return
	}
	defer src.Close()
	
	fi, err := file.Stat()
	if !assert.Nil(t, err) {
		return
	}
	size := fi.Size()
	
	head := hash.CalcHashHead(size)
	head.Hashs = slices.Collect(hash.CalcFileSubHash(head, size, file.File()))
	
	// In actual, matchIter data should be transmitted back to the client
	matchIter, err := HashMatch(context.Background(), head, src)
	if !assert.Nil(t, err) {
		return
	}
	
	// In actual, The client received the chunks transmitted by the server's matchIter
	// and concatenated them into a complete file based on the chunks
	err = SyncFile(matchIter, file)
	if !assert.Nil(t, err) {
		return
	}
	
	// md5Hash := md5.New()
	// err = SyncFileToWriter(matchIter, file, md5Hash)
	// if !assert.Nil(t, err) {
	// 	return
	// }
	
	// md5Hash = md5.New()
	// sf, _ := os.Open(source)
	// io.Copy(md5Hash, sf)
	
	// targetHash := md5Hash.Sum(nil)
	// finalHash := md5Hash.Sum(nil)
	// assert.Equal(t, targetHash, finalHash)
}
