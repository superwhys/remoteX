package match

import (
	"context"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/internal/filesystem"
)

func SyncFileToWriter(ctx context.Context, matchIter iter.Seq2[*pb.FileChunk, error], target *filesystem.File, writer io.Writer) error {
	for matchChunk, err := range matchIter {
		if err != nil {
			return errors.Wrap(err, "iter file match")
		}

		var data []byte
		if matchChunk.GetHash() != nil {
			offset := matchChunk.Hash.GetOffset()
			data, err = filesystem.BasicFs.ReadFileAtOffset(target, offset, matchChunk.Hash.GetLen())
			if err != nil {
				return errors.Wrap(err, "read file at offset")
			}
		} else {
			data = matchChunk.GetData()
		}

		_, err := writer.Write(data)
		if err != nil {
			return errors.Wrap(err, "write to Writer")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}

func SyncFile(ctx context.Context, matchIter iter.Seq2[*pb.FileChunk, error], target *filesystem.File) (err error) {
	path := target.Name()
	baseName := filepath.Base(path)
	tmpPath := filepath.Join(filepath.Dir(path), fmt.Sprintf("tmp-%s", baseName))
	tmpFile, err := filesystem.BasicFs.CreateFile(tmpPath)
	if err != nil {
		return errors.Wrap(err, "create tmp file")
	}
	defer func() {
		if err != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
			return
		}
		tmpFile.Close()
		target.Close()

		if err := os.Rename(tmpPath, path); err != nil {
			plog.Errorf("rename tmp file to target file: %v", err)
			return
		}

		f, err := filesystem.BasicFs.OpenFile(path)
		if err != nil {
			plog.Errorf("open target file after rename: %v", err)
			return
		}
		target.Update(f)
	}()

	err = SyncFileToWriter(ctx, matchIter, target, tmpFile)
	if err != nil {
		return errors.Wrap(err, "sync file to Writer")
	}

	return nil
}
