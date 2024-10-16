// File:		match_test.go
// Created by:	Hoven
// Created on:	2024-10-12
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package filesync

import (
	"context"
	"crypto/md5"
	io "io"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
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
	file, err := os.Open(target)
	if !assert.Nil(t, err) {
		return
	}

	fi, err := file.Stat()
	if !assert.Nil(t, err) {
		return
	}
	size := fi.Size()

	head := calcHashHead(size)
	head.Hashs = slices.Collect(calcFileSubHash(head, size, file))

	matchIter, err := HashMatch(context.Background(), head, "", &FileBase{Path: source})
	if !assert.Nil(t, err) {
		return
	}

	md5Hash := md5.New()
	err = SyncFile(matchIter, file, md5Hash)
	if !assert.Nil(t, err) {
		return
	}

	md5Hash = md5.New()
	sf, _ := os.Open(source)
	io.Copy(md5Hash, sf)

	targetHash := md5Hash.Sum(nil)
	finalHash := md5Hash.Sum(nil)
	assert.Equal(t, targetHash, finalHash)
}
