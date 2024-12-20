package filesync_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/plog/level"
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesystem"
	"github.com/superwhys/remoteX/pkg/protoutils"
	"golang.org/x/sync/errgroup"
)

type readWriter struct {
	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter
}

func (rw *readWriter) Close() error {
	var err error
	err = rw.ProtoMessageWriter.Close()
	err = rw.ProtoMessageReader.Close()
	if err != nil {
		return err
	}
	return nil
}

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

func TestTransferWholeFile(t *testing.T) {
	plog.Enable(level.LevelDebug)

	reader, writer := io.Pipe()
	rd := protoutils.NewProtoReader(reader)
	wr := protoutils.NewProtoWriter(writer)
	rw := &readWriter{rd, wr}

	tmpFile, err := tempFileCreate(1024 * 1024 * 100)
	if err != nil {
		t.Fatal(err)
	}
	fileName := filepath.Base(tmpFile)
	targetDest := filepath.Join(filepath.Dir(tmpFile), "sync")
	defer func() {
		os.Remove(tmpFile)
		os.Remove(targetDest)
	}()

	fmt.Println(tmpFile, targetDest)

	fileSyncer := filesync.NewFileSyncer()

	grp, ctx := errgroup.WithContext(context.Background())
	grp.Go(func() error {
		defer rw.Close()
		_, err := fileSyncer.SendFiles(ctx, rw, tmpFile, nil)
		return err
	})
	grp.Go(func() error {
		defer rw.Close()
		_, err := fileSyncer.ReceiveFiles(ctx, rw, targetDest, nil)
		return err
	})

	if err := grp.Wait(); err != nil {
		t.Fatal(err)
	}

	fs := filesystem.NewBasicFileSystem()

	tempF, err := fs.OpenFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	targetFile := filepath.Join(targetDest, fileName)
	targetF, err := fs.OpenFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	tempM, err := tempF.MD5()
	if err != nil {
		t.Fatal(err)
	}
	targetM, err := targetF.MD5()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, tempM, targetM)
}

func TestTransferDir(t *testing.T) {
	plog.Enable(level.LevelDebug)
	fileSyncer := filesync.NewFileSyncer()

	reader, writer := io.Pipe()
	rd := protoutils.NewProtoReader(reader)
	wr := protoutils.NewProtoWriter(writer)
	rw := &readWriter{rd, wr}

	source := "/Users/yong/programes/go/src/github.com/superwhys/remoteX/content"
	target := "/tmp/remoteX/content"
	grp, ctx := errgroup.WithContext(context.Background())
	grp.Go(func() error {
		defer rw.Close()
		_, err := fileSyncer.SendFiles(ctx, rw, source, nil)
		return err
	})
	grp.Go(func() error {
		defer rw.Close()
		_, err := fileSyncer.ReceiveFiles(ctx, rw, target, nil)
		return err
	})

	if err := grp.Wait(); err != nil {
		t.Fatal(err)
	}
}
