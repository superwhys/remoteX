package match

import (
	"context"
	"crypto/md5"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/remoteX/internal/filesync/hash"
	"github.com/superwhys/remoteX/internal/filesystem"
)

func tempFileCreate(size int) (string, error) {
	tmpFile, err := os.CreateTemp("/tmp/remoteX", "testfile-tmp-*.txt")
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

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "openFile(%s)", src)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func randomAddBytes(f string, fileSize int64) error {
	file, err := os.OpenFile(f, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	for range 10 {
		insertPosition := rand.Int63n(fileSize + 1)

		newData := []byte("1101010101")
		_, err = file.Seek(insertPosition, io.SeekStart)
		if err != nil {
			return err
		}

		_, err = file.Write(newData)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestFileHashMatch(t *testing.T) {
	fileSize := 1024 * 1024 * 50
	source, err := tempFileCreate(fileSize)
	if !assert.Nil(t, err) {
		return
	}

	fileName := filepath.Base(source)
	target := filepath.Join(filepath.Dir(source), "sync"+fileName)
	err = copyFile(source, target)
	if !assert.Nil(t, err) {
		return
	}

	err = randomAddBytes(source, int64(fileSize))
	err = randomAddBytes(target, int64(fileSize))
	if !assert.Nil(t, err) {
		return
	}
	defer func() {
		_ = os.Remove(target)
		_ = os.Remove(source)
	}()

	targetFile, err := filesystem.NewBasicFileSystem().OpenFile(target)
	if !assert.Nil(t, err) {
		return
	}

	src, err := filesystem.NewBasicFileSystem().OpenFile(source)
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
	t.Logf("head: %+v", head)
	head.Hashs = slices.Collect(hash.CalcFileSubHash(head, size, targetFile.File()))

	// In actual, matchIter data should be transmitted back to the client
	matchIter, err := HashMatch(context.Background(), head, src)
	if !assert.Nil(t, err) {
		return
	}

	// In actual, The client received the chunks transmitted by the server's matchIter
	// and concatenated them into a complete targetFile based on the chunks
	md5Hash := md5.New()

	for fileChunk, err := range matchIter {
		if !assert.Nil(t, err) {
			return
		}

		var data []byte
		if fileChunk.GetHash() != nil {
			offset := fileChunk.Hash.GetOffset()
			data, err = targetFile.ReadFileAtOffset(offset, fileChunk.Hash.GetLen())
			if !assert.Nil(t, err) {
				return
			}
		} else {
			data = fileChunk.GetData()
		}

		_, err := md5Hash.Write(data)
		if !assert.Nil(t, err) {
			return
		}
	}

	srcMd5, err := src.MD5()
	if !assert.Nil(t, err) {
		return
	}

	finalHash := md5Hash.Sum(nil)
	assert.Equal(t, srcMd5, finalHash)
}
