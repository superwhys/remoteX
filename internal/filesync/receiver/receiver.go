package receiver

import (
	"context"
	"iter"
	"os"
	"path/filepath"

	"github.com/go-puzzles/puzzles/putils"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/common"
	"github.com/superwhys/remoteX/internal/filesync/hash"
	"github.com/superwhys/remoteX/internal/filesync/match"
	"github.com/superwhys/remoteX/internal/filesync/opts"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/internal/filesystem"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type ReceiveTransfer struct {
	Opts      *opts.SyncOpt
	Dest      string
	DestIsDir bool
	Rw        protoutils.ProtoMessageReadWriter
}

func (rt *ReceiveTransfer) ReceiveFileList(ctx context.Context) (*pb.FileList, error) {
	var fileList pb.FileList
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		f := &pb.FileBase{}
		if err := rt.Rw.ReadMessage(f); err != nil {
			return nil, err
		}

		if f.IsEnd {
			break
		}

		if f.Entry.Type == filesystem.EntryTypeDir {
			targetDir := filepath.Join(rt.Dest, f.GetEntry().GetWpath())
			if err := rt.checkdir(targetDir); err != nil {
				return nil, errors.Wrapf(err, "checkdir: %s", targetDir)
			}
		}

		fileList.Files = append(fileList.Files, f)
		if f.GetEntry().GetRegular() {
			fileList.TotalSize += f.GetEntry().GetSize()
		}
	}

	fileList.Sort()

	return &fileList, nil
}

func (rt *ReceiveTransfer) checkdir(dir string) error {
	if putils.FileExists(dir) {
		return nil
	}

	return os.MkdirAll(dir, os.FileMode(0755))
}

func (rt *ReceiveTransfer) CalcFileHashAndSend(ctx context.Context, local string) error {
	if !putils.FileExists(local) {
		return rt.Rw.WriteMessage(&pb.HashHead{
			BlockLength: int64(common.BlockSize),
		})
	}

	in, err := filesystem.BasicFs.OpenFile(local)
	if err != nil {
		return errors.Wrapf(err, "openFile: %s", local)
	}

	fileLen, err := putils.FileSize(local)
	if err != nil {
		return errors.Wrapf(err, "calcFileSize: %s", local)
	}

	head := hash.CalcHashHead(fileLen)
	if err := rt.Rw.WriteMessage(head); err != nil {
		return err
	}

	for hb := range hash.CalcFileSubHash(head, fileLen, in.File()) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := rt.Rw.WriteMessage(hb); err != nil {
			return errors.Wrap(err, "sendHashBuf")
		}
	}

	return nil
}

func (rt *ReceiveTransfer) receiveFileChunkIter(ctx context.Context) (matchIter iter.Seq2[*pb.FileChunk, error]) {
	return func(yield func(*pb.FileChunk, error) bool) {
		for {
			var fileChunk pb.FileChunk
			if err := rt.Rw.ReadMessage(&fileChunk); err != nil {
				yield(nil, errors.Wrapf(err, "readFileChunk"))
				return
			}

			if fileChunk.GetIsEnd() || !yield(&fileChunk, nil) {
				return
			}

			select {
			case <-ctx.Done():
				yield(nil, ctx.Err())
				return
			default:
			}
		}
	}
}

func (rt *ReceiveTransfer) TransferFile(ctx context.Context, targetPath string) (err error) {
	var target *filesystem.File
	if !putils.FileExists(targetPath) {
		target, err = filesystem.BasicFs.CreateFile(targetPath)
	} else {
		target, err = filesystem.BasicFs.OpenFile(targetPath)
	}
	if err != nil {
		return errors.Wrap(err, "openFile")
	}

	matchIter := rt.receiveFileChunkIter(ctx)
	return match.SyncFile(ctx, matchIter, target)
}
