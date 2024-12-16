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
	Fs        filesystem.FileSystem
	Opts      *pb.SyncOpts
	Dest      string
	DestIsDir bool
	Rw        protoutils.ProtoMessageReadWriter
}

func (rt *ReceiveTransfer) receiveOpts() (*pb.SyncOpts, error) {
	opts := new(pb.SyncOpts)
	if err := rt.Rw.ReadMessage(opts); err != nil {
		return nil, err
	}

	return opts, nil
}

func (rt *ReceiveTransfer) MergeRemoteOpts() error {
	opts, err := rt.receiveOpts()
	if err != nil {
		return errors.Wrap(err, "receiveRemoteOpts")
	}

	rt.Opts.DryRun = opts.DryRun
	rt.Opts.Whole = opts.Whole

	return nil
}

func (rt *ReceiveTransfer) CheckDesk(fileCnt int, dest string) error {
	info, err := os.Stat(dest)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "stat")
	}

	if info == nil {
		// dest not exists
		// if multiple files are received, dest must be a folder
		// if just has one file and dest is not exists, it will be treated as a file
		if fileCnt > 1 {
			rt.DestIsDir = true
		}
	} else {
		rt.DestIsDir = info.IsDir()
	}

	if fileCnt > 1 && !rt.DestIsDir {
		return errors.New("dest is not a directory")
	}

	if rt.DestIsDir && info == nil {
		err = os.MkdirAll(dest, 0755)
	} else if !rt.DestIsDir {
		destBaseDir := filepath.Dir(dest)
		if !putils.FileExists(destBaseDir) {
			err = os.MkdirAll(destBaseDir, 0755)
		}
	}

	if err != nil {
		return errors.Wrap(err, "mkdir")
	}

	return nil
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
			if err := rt.CheckDir(targetDir); err != nil {
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

func (rt *ReceiveTransfer) CheckDir(dir string) error {
	if putils.FileExists(dir) {
		return nil
	}

	return os.MkdirAll(dir, os.FileMode(0755))
}

func (rt *ReceiveTransfer) CalcFileHashAndSend(ctx context.Context, local string) error {
	// whole file
	if rt.Opts.Whole || !putils.FileExists(local) {
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

func (rt *ReceiveTransfer) TransferFile(ctx context.Context, targetPath string) (err error) {
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
