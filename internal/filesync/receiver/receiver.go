package receiver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync"
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
	// receive filelist
	fileList, err := rt.receiveFileList(ctx)
	if err != nil {
		return errors.Wrap(err, "receiveFileList")
	}

	for idx, f := range fileList.Files {
		if err := rt.rw.WriteMessage(&filesync.FileIdx{Idx: int64(idx)}); err != nil {
			return errors.Wrapf(err, "sendFileIdx: %s", f.Name())
		}

		if opts.DryRun {
			continue
		}

		// generate each file sum and send
		fileExists, err := rt.calcFileHashAndSend(f)
		if err != nil {
			return errors.Wrapf(err, "calculate file hash: %s", f.GetPath())
		}

		// receive server file match chunk
		if err := rt.transferFile(fileExists, dest, f); err != nil {
			return errors.Wrap(err, "transferFile")
		}
	}

	return nil
}
