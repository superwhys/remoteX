package sender

import (
	"context"
	"path/filepath"
	"strings"
	
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesync/opts"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type FileSender interface {
	SendFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, path string, opts *opts.SyncOpt) (err error)
}

type Sender struct{}

func (s *Sender) SendFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, path string, opts *opts.SyncOpt) (err error) {
	
	st := &sendTransfer{
		opts: opts,
		rw:   rw,
	}
	// TODO: 1. send whole file list
	fileList, err := st.sendFileList(ctx, path)
	if err != nil {
		return errors.Wrap(err, "sendFileList")
	}
	_ = fileList
	// 2. receive client process file idx
	// 3. receive client file sums
	
	return nil
}

type sendTransfer struct {
	opts *opts.SyncOpt
	rw   protoutils.ProtoMessageReadWriter
}

func (st *sendTransfer) sendFileList(ctx context.Context, root string) (*filesync.FileList, error) {
	var fileList filesync.FileList
	
	strip := filepath.Dir(filepath.Clean(root)) + "/"
	if strings.HasSuffix(root, "/") {
		strip = filepath.Clean(root) + "/"
	}
	
	fileList.Strip = strip
	
	for f := range filesync.GetFileList(root) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		fileList.Files = append(fileList.Files, f)
		fileList.TotalSize += f.GetSize()
		
		if err := st.rw.WriteMessage(f); err != nil {
			return nil, errors.Wrapf(err, "write file: %s", f.GetPath())
		}
	}
	
	fileList.Sort()
	
	if err := st.rw.WriteMessage(&filesync.FileBase{IsEnd: true}); err != nil {
		return nil, errors.Wrap(err, "write end")
	}
	
	return &fileList, nil
}
