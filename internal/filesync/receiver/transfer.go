package receiver

import (
	"context"
	"iter"
	"os"
	"path/filepath"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesync/file"
	"github.com/superwhys/remoteX/internal/filesync/hash"
	"github.com/superwhys/remoteX/internal/filesync/match"
	"github.com/superwhys/remoteX/internal/filesync/opts"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type receiveTransfer struct {
	opts *opts.SyncOpt
	dest string
	rw   protoutils.ProtoMessageReadWriter
}

func (rt *receiveTransfer) receiveFileList(ctx context.Context) (*filesync.FileList, error) {
	var fileList filesync.FileList
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		f := &filesync.FileBase{}
		if err := rt.rw.ReadMessage(f); err != nil {
			return nil, err
		}

		if f.IsEnd {
			break
		}

		fileList.Files = append(fileList.Files, f)
		if f.Regular {
			fileList.TotalSize += f.Size
		}
	}

	fileList.Sort()

	return &fileList, nil
}

func (rt *receiveTransfer) calcFileHashAndSend(f *filesync.FileBase) (exists bool, err error) {
	local := filepath.Join(rt.dest, f.Name())
	st, err := os.Stat(local)
	if os.IsNotExist(err) {
		// need whole file
		return false, rt.rw.WriteMessage(&filesync.HashHead{
			BlockLength: int64(hash.GetDefaultBlockSize()),
		})
	} else if err != nil {
		return false, err
	}

	in, err := file.OpenFile(local)
	if err != nil {
		return false, errors.Wrapf(err, "openFile: %s", local)
	}
	fileLen := st.Size()

	head := hash.CalcHashHead(fileLen)
	if err := rt.rw.WriteMessage(head); err != nil {
		return false, err
	}

	for hb := range hash.CalcFileSubHash(head, fileLen, in.File()) {
		if err := rt.rw.WriteMessage(hb); err != nil {
			return false, errors.Wrap(err, "sendHashBuf")
		}
	}

	return true, nil
}

func (rt *receiveTransfer) receiveFileChunkIter() (matchIter iter.Seq2[*filesync.FileChunk, error]) {
	return func(yield func(*filesync.FileChunk, error) bool) {
		for {
			var fileChunk filesync.FileChunk
			if err := rt.rw.ReadMessage(&fileChunk); err != nil {
				yield(nil, errors.Wrapf(err, "readFileChunk"))
				return
			}

			if fileChunk.GetIsEnd() || !yield(&fileChunk, nil) {
				return
			}
		}
	}
}

func (rt *receiveTransfer) transferFile(fileExists bool, dest string, fb *filesync.FileBase) (err error) {
	var target *file.File

	targetPath := filepath.Join(dest, fb.GetPath())
	if !fileExists {
		target, err = file.CreateFile(targetPath)
	} else {
		target, err = file.OpenFile(targetPath)
	}
	if err != nil {
		return errors.Wrapf(err, "openFile: %s", targetPath)
	}

	plog.Debugf("open target file %s success", targetPath)

	matchIter := rt.receiveFileChunkIter()
	return match.SyncFile(matchIter, target)
}
