package receiver

import (
	"context"
	"os"
	"path/filepath"
	
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesync/file"
	"github.com/superwhys/remoteX/internal/filesync/hash"
	"github.com/superwhys/remoteX/internal/filesync/opts"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type FileReceiver interface {
	ReceiveFile(ctx context.Context, rw protoutils.ProtoMessageReadWriter, dest string, opts *opts.SyncOpt) error
}

type Receiver struct{}

func (r *Receiver) ReceiveFile(ctx context.Context, rw protoutils.ProtoMessageReadWriter, dest string, opts *opts.SyncOpt) error {
	rt := &receiveTransfer{
		opts: opts,
		dest: dest,
		rw:   rw,
	}
	// TODO: 1. receive filelist
	fileList, err := rt.receiveFileList(ctx)
	if err != nil {
		return errors.Wrap(err, "receiveFileList")
	}
	_ = fileList
	// 2. generate each file sum and send
	
	// 3. receive file
	
	return nil
}

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

func (rt *receiveTransfer) generateFilesHash(ctx context.Context, fileList *filesync.FileList) error {
	for idx, f := range fileList.Files {
		if err := rt.calcFileHashAndSend(ctx, int64(idx), f); err != nil {
			return errors.Wrapf(err, "calculate file hash: %s", f.GetPath())
		}
	}
	
	return nil
}

func (rt *receiveTransfer) calcFileHashAndSend(ctx context.Context, idx int64, f *filesync.FileBase) error {
	local := filepath.Join(rt.dest, f.Name())
	st, err := os.Stat(local)
	if os.IsNotExist(err) {
		// TODO: need whole file
		
	} else if err != nil {
		return err
	}
	
	in, err := file.OpenFile(local)
	if err != nil {
		return errors.Wrapf(err, "openFile: %s", local)
	}
	fileLen := st.Size()
	
	// 1. write idx
	if err := rt.rw.WriteMessage(&filesync.FileIdx{Idx: idx}); err != nil {
		return errors.Wrapf(err, "sendFileIdx: %s", local)
	}
	
	// 2. calc hash
	head := hash.CalcHashHead(fileLen)
	if err := rt.rw.WriteMessage(head); err != nil {
		return err
	}
	
	for hb := range hash.CalcFileSubHash(head, fileLen, in.File()) {
		if err := rt.rw.WriteMessage(hb); err != nil {
			return errors.Wrap(err, "sendHashBuf")
		}
	}
	
	return nil
}
