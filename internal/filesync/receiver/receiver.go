package receiver

import (
	"context"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/common"
	"github.com/superwhys/remoteX/internal/filesync/hash"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/internal/filesystem"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type ReceiveTransfer struct {
	Fs            filesystem.FileSystem
	Rw            protoutils.ProtoMessageReadWriter
	Dest          string
	ActualReceive int
}

func (rt *ReceiveTransfer) TransferSize() int64 {
	return int64(rt.ActualReceive)
}

func (rt *ReceiveTransfer) Transfer(ctx context.Context, file *pb.FileBase, opts *pb.SyncOpts) error {
	targetPath := filepath.Join(rt.Dest, file.GetEntry().GetWpath())
	ctx = plog.With(ctx, "file", targetPath)

	// generate each file sum and send
	plog.Debugc(ctx, "start calc file hash")
	err := rt.calcFileHashAndSend(ctx, targetPath, opts)
	if err != nil {
		return errors.Wrapf(err, "calculate file hash: %s", targetPath)
	}
	plog.Debugc(ctx, "calc file hash success")

	// receive server file match chunk
	plog.Debugc(ctx, "start transfer file")
	if err := rt.receiveFile(ctx, targetPath); err != nil {
		return errors.Wrap(err, "transferFile")
	}
	return nil
}

func (rt *ReceiveTransfer) calcFileHashAndSend(ctx context.Context, local string, opts *pb.SyncOpts) error {
	// whole file
	if opts.Whole || !putils.FileExists(local) {
		return rt.Rw.WriteMessage(&pb.HashHead{
			BlockLength: int64(common.BlockSize),
		})
	}

	in, err := rt.Fs.OpenFile(local)
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

func (rt *ReceiveTransfer) receiveFile(ctx context.Context, targetPath string) (err error) {
	var target filesystem.File
	if !putils.FileExists(targetPath) {
		target, err = rt.Fs.CreateFile(targetPath)
	} else {
		target, err = rt.Fs.OpenFile(targetPath)
	}
	if err != nil {
		return errors.Wrap(err, "openFile")
	}

	fileChunkIter := rt.receiveFileChunkIter(ctx)
	return rt.syncFile(ctx, fileChunkIter, target)
}

func (rt *ReceiveTransfer) receiveFileChunkIter(ctx context.Context) (fileChunkIter iter.Seq2[*pb.FileChunk, error]) {
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

func (rt *ReceiveTransfer) syncFileToWriter(ctx context.Context, fileChunkIter iter.Seq2[*pb.FileChunk, error], target filesystem.File, writer io.Writer) error {
	for fileChunk, err := range fileChunkIter {
		if err != nil {
			return errors.Wrap(err, "iter file match")
		}

		var data []byte
		if fileChunk.GetHash() != nil {
			offset := fileChunk.Hash.GetOffset()
			data, err = target.ReadFileAtOffset(offset, fileChunk.Hash.GetLen())
			if err != nil {
				return errors.Wrap(err, "read file at offset")
			}
		} else {
			data = fileChunk.GetData()
		}

		_, err := writer.Write(data)
		if err != nil {
			return errors.Wrap(err, "write to Writer")
		}

		rt.ActualReceive += len(fileChunk.GetData())

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}

func (rt *ReceiveTransfer) syncFile(ctx context.Context, fileChunkIter iter.Seq2[*pb.FileChunk, error], target filesystem.File) (err error) {
	path := target.Name()
	baseName := filepath.Base(path)
	tmpPath := filepath.Join(filepath.Dir(path), fmt.Sprintf("tmp-%s", baseName))
	tmpFile, err := rt.Fs.CreateFile(tmpPath)
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

		f, err := rt.Fs.OpenFile(path)
		if err != nil {
			plog.Errorf("open target file after rename: %v", err)
			return
		}
		target.Update(f)
	}()

	err = rt.syncFileToWriter(ctx, fileChunkIter, target, tmpFile)
	if err != nil {
		return errors.Wrap(err, "sync file to Writer")
	}

	return nil
}
